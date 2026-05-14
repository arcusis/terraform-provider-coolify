package resources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// coolify_envs_bulk manages all environment variables for an application, service,
// or database as a single atomic set via the bulk API endpoint.

type envsBulkModel struct {
	ID           types.String `tfsdk:"id"`
	ResourceType types.String `tfsdk:"resource_type"`
	ResourceUUID types.String `tfsdk:"resource_uuid"`
	Variables    types.Map    `tfsdk:"variables"`
}

type envsBulkResource struct{ client *coolify.Client }

func NewEnvsBulkResource() resource.Resource { return &envsBulkResource{} }

func (r *envsBulkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_envs_bulk"
}

func (r *envsBulkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"resource_type": schema.StringAttribute{Required: true, Description: "Type: application, service, or database."},
			"resource_uuid": schema.StringAttribute{Required: true},
			"variables":     schema.MapAttribute{Required: true, ElementType: types.StringType},
		},
	}
}

func (r *envsBulkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		r.client = client
	}
}

func (r *envsBulkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan envsBulkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.sendBulk(ctx, plan); err != nil {
		resp.Diagnostics.AddError("Unable to bulk-set environment variables", err.Error())
		return
	}

	plan.ID = types.StringValue(plan.ResourceType.ValueString() + "/" + plan.ResourceUUID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *envsBulkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state envsBulkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read current env vars and update state
	var out any
	envsPath := fmt.Sprintf("/api/v1/%ss/%s/envs", state.ResourceType.ValueString(), state.ResourceUUID.ValueString())
	if err := r.client.Get(ctx, envsPath, &out); err == nil {
		vars := map[string]attr.Value{}
		for _, item := range envList(out) {
			k := stringFromAny(item["key"])
			v := stringFromAny(item["value"])
			if k != "" {
				vars[k] = types.StringValue(v)
			}
		}
		m, diags := types.MapValue(types.StringType, vars)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.Variables = m
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *envsBulkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan envsBulkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.sendBulk(ctx, plan); err != nil {
		resp.Diagnostics.AddError("Unable to bulk-update environment variables", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *envsBulkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Bulk delete by sending empty array
	var state envsBulkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bulkPath := fmt.Sprintf("/api/v1/%ss/%s/envs/bulk", state.ResourceType.ValueString(), state.ResourceUUID.ValueString())
	_ = r.client.Patch(ctx, bulkPath, map[string]any{"data": []any{}}, nil)
}

func (r *envsBulkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *envsBulkResource) sendBulk(ctx context.Context, model envsBulkModel) error {
	items := make([]map[string]any, 0, len(model.Variables.Elements()))
	for k, v := range model.Variables.Elements() {
		sv, ok := v.(types.String)
		if !ok {
			continue
		}
		items = append(items, map[string]any{"key": k, "value": sv.ValueString()})
	}
	bulkPath := fmt.Sprintf("/api/v1/%ss/%s/envs/bulk", model.ResourceType.ValueString(), model.ResourceUUID.ValueString())
	return r.client.Patch(ctx, bulkPath, map[string]any{"data": items}, nil)
}
