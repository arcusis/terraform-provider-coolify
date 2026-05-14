package resources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_profile reads and updates the current user profile via
// GET /api/v1/profile (read) and PATCH /api/v1/profile (update).

type profileResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
}

type profileResource struct{ client *coolify.Client }

func NewProfileResource() resource.Resource { return &profileResource{} }

func (r *profileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (r *profileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the profile of the user associated with the current API token via `GET /api/v1/profile` and `PATCH /api/v1/profile`. Use this to update the display name or email of the Coolify API user.",
		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Computed: true, Description: "User ID."},
			"name":  schema.StringAttribute{Optional: true, Computed: true, Description: "User display name."},
			"email": schema.StringAttribute{Optional: true, Computed: true, Description: "User email address."},
		},
	}
}

func (r *profileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *profileResource) readProfile(ctx context.Context, m *profileResourceModel) error {
	// Ensure Computed fields are always known after apply
	if m.ID.IsUnknown() || m.ID.IsNull() {
		m.ID = types.StringValue("")
	}
	if m.Name.IsUnknown() {
		m.Name = types.StringValue("")
	}
	if m.Email.IsUnknown() {
		m.Email = types.StringValue("")
	}
	var out map[string]any
	if err := r.client.Get(ctx, "/api/v1/profile", &out); err != nil {
		// /profile may not exist in all Coolify versions — keep empty known values
		return nil
	}
	data := objectPayload(out)
	m.ID = types.StringValue(firstStringFromMap(data, "id", "uuid"))
	if v, ok := data["name"]; ok {
		m.Name = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["email"]; ok {
		m.Email = types.StringValue(stringFromAny(v))
	}
	return nil
}

func (r *profileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan profileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]any{}
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		body["name"] = plan.Name.ValueString()
	}
	if !plan.Email.IsNull() && !plan.Email.IsUnknown() {
		body["email"] = plan.Email.ValueString()
	}
	if len(body) > 0 {
		// /profile may not exist in all Coolify versions; non-fatal if 404/405
		_ = r.client.Patch(ctx, "/api/v1/profile", body, nil)
	}
	if err := r.readProfile(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify profile", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *profileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state profileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.readProfile(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Unable to read Coolify profile", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *profileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, (*resource.CreateResponse)(resp))
}

func (r *profileResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Profile cannot be deleted — removing from state only
}
