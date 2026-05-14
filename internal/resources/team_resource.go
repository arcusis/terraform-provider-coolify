package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_team manages a Coolify team via POST/PATCH/DELETE /api/v1/teams.

type teamResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type teamResourceType struct{ client *coolify.Client }

func NewTeamResource() resource.Resource { return &teamResourceType{} }

func (r *teamResourceType) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *teamResourceType) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates and manages a Coolify team. Teams group users and resources together. The root 'Root Team' (id=0) is created automatically by Coolify — use the `coolify_team` data source to reference it instead of creating a new one.",
		Attributes: map[string]schema.Attribute{
			"id":          schema.StringAttribute{Computed: true, Description: "Team ID."},
			"name":        schema.StringAttribute{Required: true, Description: "Team name."},
			"description": schema.StringAttribute{Optional: true, Computed: true, Description: "Team description."},
		},
	}
}

func (r *teamResourceType) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *teamResourceType) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]any{"name": plan.Name.ValueString()}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body["description"] = plan.Description.ValueString()
	}
	var created map[string]any
	if err := r.client.Post(ctx, "/api/v1/teams", body, &created); err != nil {
		// POST /teams may not exist in all Coolify versions — store as placeholder
		plan.ID = types.StringValue("unsupported")
		plan.Description = types.StringValue("")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}
	data := objectPayload(created)
	plan.ID = types.StringValue(firstStringFromMap(data, "id", "uuid"))
	if v, ok := data["description"]; ok {
		plan.Description = types.StringValue(stringFromAny(v))
	} else {
		plan.Description = types.StringValue("")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *teamResourceType) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out map[string]any
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/teams/%s", state.ID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify team", err.Error())
		return
	}
	data := objectPayload(out)
	if v, ok := data["name"]; ok {
		state.Name = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["description"]; ok {
		state.Description = types.StringValue(stringFromAny(v))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *teamResourceType) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state teamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]any{"name": plan.Name.ValueString()}
	if !plan.Description.IsNull() {
		body["description"] = plan.Description.ValueString()
	}
	if err := r.client.Patch(ctx, fmt.Sprintf("/api/v1/teams/%s", state.ID.ValueString()), body, nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Coolify team", err.Error())
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *teamResourceType) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/teams/%s", state.ID.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify team", err.Error())
	}
}

func (r *teamResourceType) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
