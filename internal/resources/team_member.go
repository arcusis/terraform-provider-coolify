package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_team_member manages team membership via
// POST /api/v1/teams/{id}/members (add) and
// DELETE /api/v1/teams/{id}/members/{user_id} (remove).

type teamMemberResourceModel struct {
	TeamID types.String `tfsdk:"team_id"`
	UserID types.String `tfsdk:"user_id"`
	Email  types.String `tfsdk:"email"`
	Role   types.String `tfsdk:"role"`
}

type teamMemberResource struct{ client *coolify.Client }

func NewTeamMemberResource() resource.Resource { return &teamMemberResource{} }

func (r *teamMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_member"
}

func (r *teamMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Adds a user to a Coolify team via `POST /api/v1/teams/{id}/members`. Destroying this resource removes the user from the team via `DELETE /api/v1/teams/{id}/members/{user_id}`.",
		Attributes: map[string]schema.Attribute{
			"team_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the team to add the user to.",
			},
			"user_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the user to add. If omitted, `email` is used to look up the user.",
			},
			"email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Email address of the user to add.",
			},
			"role": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Role to assign within the team (`owner` or `member`).",
			},
		},
	}
}

func (r *teamMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *teamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]any{}
	if !plan.UserID.IsNull() && !plan.UserID.IsUnknown() && plan.UserID.ValueString() != "" {
		body["user_id"] = plan.UserID.ValueString()
	}
	if !plan.Email.IsNull() && !plan.Email.IsUnknown() && plan.Email.ValueString() != "" {
		body["email"] = plan.Email.ValueString()
	}
	if !plan.Role.IsNull() && !plan.Role.IsUnknown() {
		body["role"] = plan.Role.ValueString()
	}

	var created map[string]any
	path := fmt.Sprintf("/api/v1/teams/%s/members", plan.TeamID.ValueString())
	if err := r.client.Post(ctx, path, body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to add member to Coolify team", err.Error())
		return
	}
	data := objectPayload(created)
	if uid := firstStringFromMap(data, "id", "user_id", "uuid"); uid != "" {
		plan.UserID = types.StringValue(uid)
	}
	if v, ok := data["email"]; ok {
		plan.Email = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["role"]; ok {
		plan.Role = types.StringValue(stringFromAny(v))
	}
	if plan.Role.IsNull() || plan.Role.IsUnknown() {
		plan.Role = types.StringValue("member")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *teamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Verify membership still exists by listing team members
	var out any
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/teams/%s/members", state.TeamID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		// Non-fatal read error — keep state
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	found := false
	if items, ok := out.([]any); ok {
		for _, raw := range items {
			if item, ok := raw.(map[string]any); ok {
				id := firstStringFromMap(item, "id", "uuid")
				if id == state.UserID.ValueString() {
					found = true
					if v, ok := item["email"]; ok {
						state.Email = types.StringValue(stringFromAny(v))
					}
					break
				}
			}
		}
	} else {
		found = true // can't verify, keep state
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *teamMemberResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *teamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	path := fmt.Sprintf("/api/v1/teams/%s/members/%s", state.TeamID.ValueString(), state.UserID.ValueString())
	if err := r.client.Delete(ctx, path, nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to remove member from Coolify team", err.Error())
	}
}
