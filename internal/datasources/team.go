package datasources

import (
	"context"
	"fmt"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type teamDataSource struct {
	client *coolify.Client
}

type teamDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Coolify team ID.",
			},
			"name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Coolify team name. Used for lookup when id is omitted.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Coolify team description.",
			},
		},
	}
}

func (d *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*coolify.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *coolify.Client from provider configuration.")
		return
	}
	d.client = client
}

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config teamDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data map[string]any
	if !config.ID.IsNull() && !config.ID.IsUnknown() && config.ID.ValueString() != "" {
		var out map[string]any
		if err := d.client.Get(ctx, fmt.Sprintf("/api/v1/teams/%s", config.ID.ValueString()), &out); err != nil {
			resp.Diagnostics.AddError("Unable to read Coolify team data source", err.Error())
			return
		}
		data = objectPayload(out)
	} else {
		var out any
		if err := d.client.Get(ctx, "/api/v1/teams", &out); err != nil {
			resp.Diagnostics.AddError("Unable to list Coolify teams", err.Error())
			return
		}
		data = findTeam(out, config.Name.ValueString())
		if len(data) == 0 {
			resp.Diagnostics.AddError("Coolify team not found", "No Coolify team matched the configured name.")
			return
		}
	}

	if id := firstString(data, "id", "uuid"); id != "" {
		config.ID = types.StringValue(id)
	}
	if value, ok := data["name"]; ok {
		config.Name = types.StringValue(stringFromAny(value))
	}
	if value, ok := data["description"]; ok {
		config.Description = types.StringValue(stringFromAny(value))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

func findTeam(out any, name string) map[string]any {
	for _, item := range dataList(out) {
		if name == "" || stringFromAny(item["name"]) == name {
			return item
		}
	}
	return map[string]any{}
}

func dataList(out any) []map[string]any {
	switch v := out.(type) {
	case []any:
		return mapList(v)
	case map[string]any:
		if data, ok := v["data"].([]any); ok {
			return mapList(data)
		}
		if teams, ok := v["teams"].([]any); ok {
			return mapList(teams)
		}
		return []map[string]any{objectPayload(v)}
	default:
		return nil
	}
}

func mapList(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if mapped, ok := item.(map[string]any); ok {
			result = append(result, mapped)
		}
	}
	return result
}

func objectPayload(data map[string]any) map[string]any {
	if nested, ok := data["data"].(map[string]any); ok {
		return nested
	}
	return data
}

func firstString(data map[string]any, keys ...string) string {
	data = objectPayload(data)
	for _, key := range keys {
		if value, ok := data[key]; ok {
			if text := stringFromAny(value); text != "" {
				return text
			}
		}
	}
	return ""
}

func stringFromAny(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprint(v)
	}
}

func int64FromAny(value any) int64 {
	switch v := value.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}

func boolFromAny(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	default:
		return false
	}
}
