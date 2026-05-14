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

// coolify_server_hetzner provisions a new Hetzner Cloud server and registers
// it with Coolify in a single API call (POST /api/v1/servers/hetzner).

type serverHetznerModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	CloudTokenUUID types.String `tfsdk:"cloud_token_uuid"`
	ServerType     types.String `tfsdk:"server_type"`
	Location       types.String `tfsdk:"location"`
	Image          types.String `tfsdk:"image"`
	PrivateKeyUUID types.String `tfsdk:"private_key_uuid"`
	Description    types.String `tfsdk:"description"`
	IPAddress      types.String `tfsdk:"ip_address"`
}

type serverHetznerResource struct{ client *coolify.Client }

func NewServerHetznerResource() resource.Resource { return &serverHetznerResource{} }

func (r *serverHetznerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_hetzner"
}

func (r *serverHetznerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provisions a new Hetzner Cloud server and registers it with Coolify in a single step. Requires a `coolify_cloud_token` with Hetzner credentials. Destroying this resource deletes the server registration from Coolify but does NOT terminate the Hetzner server — manage the Hetzner VM lifecycle separately.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Coolify server UUID assigned after creation.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name for the server in both Coolify and Hetzner.",
			},
			"cloud_token_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the `coolify_cloud_token` containing Hetzner credentials.",
			},
			"server_type": schema.StringAttribute{
				Required:    true,
				Description: "Hetzner server type name (e.g. `cpx11`, `cpx21`, `cx22`). Use `coolify_hetzner_server_types` to list available types.",
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "Hetzner datacenter location name (e.g. `nbg1`, `fsn1`, `hel1`). Use `coolify_hetzner_locations` to list available locations.",
			},
			"image": schema.StringAttribute{
				Required:    true,
				Description: "Hetzner OS image name (e.g. `ubuntu-22.04`). Use `coolify_hetzner_images` to list available images.",
			},
			"private_key_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the `coolify_private_key` to install on the new server for SSH access.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Optional description for the server.",
			},
			"ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "Public IP address of the provisioned Hetzner server.",
			},
		},
	}
}

func (r *serverHetznerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *serverHetznerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverHetznerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := map[string]any{
		"name":             plan.Name.ValueString(),
		"cloud_token_uuid": plan.CloudTokenUUID.ValueString(),
		"server_type":      plan.ServerType.ValueString(),
		"location":         plan.Location.ValueString(),
		"image":            plan.Image.ValueString(),
		"private_key_uuid": plan.PrivateKeyUUID.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		body["description"] = plan.Description.ValueString()
	}

	var created map[string]any
	if err := r.client.Post(ctx, "/api/v1/servers/hetzner", body, &created); err != nil {
		resp.Diagnostics.AddError("Unable to create Hetzner server", err.Error())
		return
	}

	id := firstStringFromMap(created, "uuid", "id")
	if id == "" {
		resp.Diagnostics.AddError("Create response missing ID", "Hetzner server create did not return an ID.")
		return
	}

	plan.ID = types.StringValue(id)
	if v, ok := created["ip"]; ok {
		plan.IPAddress = types.StringValue(stringFromAny(v))
	} else {
		plan.IPAddress = types.StringValue("")
	}
	if plan.Description.IsNull() || plan.Description.IsUnknown() {
		plan.Description = types.StringValue("")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverHetznerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverHetznerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out map[string]any
	if err := r.client.Get(ctx, fmt.Sprintf("/api/v1/servers/%s", state.ID.ValueString()), &out); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Unable to read Hetzner server", err.Error())
		return
	}
	data := objectPayload(out)
	if v, ok := data["name"]; ok {
		state.Name = types.StringValue(stringFromAny(v))
	}
	if v, ok := data["description"]; ok {
		state.Description = types.StringValue(stringFromAny(v))
	}
	ip := stringFromAny(data["ip"])
	if ip == "" {
		ip = stringFromAny(data["hostname"])
	}
	state.IPAddress = types.StringValue(ip)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serverHetznerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state serverHetznerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := map[string]any{
		"name": plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		body["description"] = plan.Description.ValueString()
	}
	if err := r.client.Patch(ctx, fmt.Sprintf("/api/v1/servers/%s", state.ID.ValueString()), body, nil); err != nil {
		resp.Diagnostics.AddError("Unable to update Hetzner server", err.Error())
		return
	}
	plan.ID = state.ID
	plan.IPAddress = state.IPAddress
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverHetznerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serverHetznerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.Delete(ctx, fmt.Sprintf("/api/v1/servers/%s", state.ID.ValueString()), nil); err != nil {
		if httpErr, ok := err.(*coolify.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Unable to delete Hetzner server", err.Error())
	}
}

func (r *serverHetznerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
