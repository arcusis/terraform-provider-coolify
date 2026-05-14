package resources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_server_validate triggers GET /servers/{uuid}/validate and surfaces
// the validation result. Deleting removes the resource from state only (no API side-effect).

type serverValidateModel struct {
	ServerUUID types.String `tfsdk:"server_uuid"`
	Status     types.String `tfsdk:"status"`
}

type serverValidateResource struct{ client *coolify.Client }

func NewServerValidateResource() resource.Resource { return &serverValidateResource{} }

func (r *serverValidateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_validate"
}

func (r *serverValidateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Validates that Coolify can connect to a server over SSH. Creating this resource triggers the validation check; the result is stored as `status`. Use it as a dependency gate—downstream resources that require a healthy server can `depends_on` this resource.",
		Attributes: map[string]schema.Attribute{
			"server_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the server to validate.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Validation result returned by Coolify (e.g. OK, unreachable).",
			},
		},
	}
}

func (r *serverValidateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *serverValidateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverValidateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	path := fmt.Sprintf("/api/v1/servers/%s/validate", plan.ServerUUID.ValueString())
	if err := r.client.Get(ctx, path, &out); err != nil {
		// Non-fatal: validation failure is surfaced as status, not as an error.
		plan.Status = types.StringValue("unreachable")
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	status := "ok"
	if m, ok := out.(map[string]any); ok {
		if v := m["status"]; v != nil {
			status = fmt.Sprint(v)
		} else if v := m["message"]; v != nil {
			status = fmt.Sprint(v)
		}
	}
	plan.Status = types.StringValue(status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverValidateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Re-validate on refresh
	var state serverValidateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	path := fmt.Sprintf("/api/v1/servers/%s/validate", state.ServerUUID.ValueString())
	if err := r.client.Get(ctx, path, &out); err != nil {
		state.Status = types.StringValue("unreachable")
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}
	status := "ok"
	if m, ok := out.(map[string]any); ok {
		if v := m["status"]; v != nil {
			status = fmt.Sprint(v)
		}
	}
	state.Status = types.StringValue(status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serverValidateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// server_uuid is required — it can't change without recreate
}

func (r *serverValidateResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Nothing to delete on the API side
}
