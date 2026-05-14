package resources

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type environmentModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectUUID types.String `tfsdk:"project_uuid"`
	Name        types.String `tfsdk:"name"`
}

type environmentResource struct {
	client *coolify.Client
}

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, Description: "Composite ID: project_uuid/environment_name."},
			"project_uuid": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
		},
	}
}

func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client.")
		return
	}
	r.client = client
}

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan environmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{"name": plan.Name.ValueString()}

	var created map[string]any
	if err := r.client.Post(ctx, fmt.Sprintf("/api/v1/projects/%s/environments", plan.ProjectUUID.ValueString()), body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Coolify environment", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.ProjectUUID.ValueString() + "/" + plan.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectUUID, envName := splitEnvID(state.ID.ValueString())
	if projectUUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	var out map[string]any
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/projects/%s/%s", projectUUID, envName), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Coolify environment", err.Error())
		return
	}

	out = objectPayload(out)
	if name, ok := out["name"]; ok {
		state.Name = types.StringValue(stringFromAny(name))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *environmentResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Environments cannot be updated", "Coolify environments do not support in-place update; destroy and recreate.")
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state environmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectUUID, envName := splitEnvID(state.ID.ValueString())
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/projects/%s/environments/%s", projectUUID, envName), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Coolify environment", err.Error())
	}
}

func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func splitEnvID(id string) (projectUUID, envName string) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
