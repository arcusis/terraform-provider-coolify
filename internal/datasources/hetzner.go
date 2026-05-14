package datasources

import (
	"context"

	"github.com/arcusis/terraform-provider-coolify/internal/coolify"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ── Hetzner Locations ────────────────────────────────────────────────────────

type hetznerLocationModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	City        types.String `tfsdk:"city"`
	Country     types.String `tfsdk:"country"`
}

type hetznerLocationsDataSourceModel struct {
	CloudTokenUUID types.String           `tfsdk:"cloud_token_uuid"`
	Locations      []hetznerLocationModel `tfsdk:"locations"`
}

type hetznerLocationsDataSource struct{ client *coolify.Client }

func NewHetznerLocationsDataSource() datasource.DataSource { return &hetznerLocationsDataSource{} }

func (d *hetznerLocationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hetzner_locations"
}

func (d *hetznerLocationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{Required: true},
			"locations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"city":        schema.StringAttribute{Computed: true},
						"country":     schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *hetznerLocationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *hetznerLocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config hetznerLocationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, "/api/v1/hetzner/locations?cloud_token_uuid="+config.CloudTokenUUID.ValueString(), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Hetzner locations", err.Error())
		return
	}
	config.Locations = []hetznerLocationModel{}
	for _, item := range dataList(out) {
		config.Locations = append(config.Locations, hetznerLocationModel{
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
			City:        types.StringValue(stringFromAny(item["city"])),
			Country:     types.StringValue(stringFromAny(item["country"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Hetzner Server Types ──────────────────────────────────────────────────────

type hetznerServerTypeModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Cores       types.Int64  `tfsdk:"cores"`
	Memory      types.Int64  `tfsdk:"memory"`
}

type hetznerServerTypesDataSourceModel struct {
	CloudTokenUUID types.String             `tfsdk:"cloud_token_uuid"`
	ServerTypes    []hetznerServerTypeModel `tfsdk:"server_types"`
}

type hetznerServerTypesDataSource struct{ client *coolify.Client }

func NewHetznerServerTypesDataSource() datasource.DataSource {
	return &hetznerServerTypesDataSource{}
}

func (d *hetznerServerTypesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hetzner_server_types"
}

func (d *hetznerServerTypesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{Required: true},
			"server_types": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"cores":       schema.Int64Attribute{Computed: true},
						"memory":      schema.Int64Attribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *hetznerServerTypesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *hetznerServerTypesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config hetznerServerTypesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var out any
	if err := d.client.Get(ctx, "/api/v1/hetzner/server-types?cloud_token_uuid="+config.CloudTokenUUID.ValueString(), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Hetzner server types", err.Error())
		return
	}
	config.ServerTypes = []hetznerServerTypeModel{}
	for _, item := range dataList(out) {
		config.ServerTypes = append(config.ServerTypes, hetznerServerTypeModel{
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
			Cores:       types.Int64Value(int64FromAny(item["cores"])),
			Memory:      types.Int64Value(int64FromAny(item["memory"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Hetzner Images ─────────────────────────────────────────────────────────────

type hetznerImageModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	OsFamily    types.String `tfsdk:"os_family"`
	OsVersion   types.String `tfsdk:"os_version"`
}

type hetznerImagesDataSourceModel struct {
	CloudTokenUUID types.String        `tfsdk:"cloud_token_uuid"`
	Images         []hetznerImageModel `tfsdk:"images"`
}

type hetznerImagesDataSource struct{ client *coolify.Client }

func NewHetznerImagesDataSource() datasource.DataSource { return &hetznerImagesDataSource{} }

func (d *hetznerImagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hetzner_images"
}

func (d *hetznerImagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists available OS images in your Hetzner Cloud account. Requires a cloud token with Hetzner credentials.",
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the Hetzner cloud token stored in Coolify (from `coolify_cloud_token`).",
			},
			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Available Hetzner OS images.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.Int64Attribute{Computed: true, Description: "Hetzner image ID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Image name (e.g. ubuntu-22.04)."},
						"description": schema.StringAttribute{Computed: true, Description: "Human-readable description."},
						"os_family":   schema.StringAttribute{Computed: true, Description: "OS family (ubuntu, debian, centos, etc.)."},
						"os_version":  schema.StringAttribute{Computed: true, Description: "OS version string."},
					},
				},
			},
		},
	}
}

func (d *hetznerImagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *hetznerImagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config hetznerImagesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/hetzner/images?cloud_token_uuid="+config.CloudTokenUUID.ValueString(), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Hetzner images", err.Error())
		return
	}
	config.Images = []hetznerImageModel{}
	for _, item := range dataList(out) {
		config.Images = append(config.Images, hetznerImageModel{
			ID:          types.Int64Value(int64FromAny(item["id"])),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Description: types.StringValue(stringFromAny(item["description"])),
			OsFamily:    types.StringValue(stringFromAny(item["os_flavor"])),
			OsVersion:   types.StringValue(stringFromAny(item["os_version"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}

// ── Hetzner SSH Keys ──────────────────────────────────────────────────────────

type hetznerSSHKeyModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Fingerprint types.String `tfsdk:"fingerprint"`
}

type hetznerSSHKeysDataSourceModel struct {
	CloudTokenUUID types.String         `tfsdk:"cloud_token_uuid"`
	SSHKeys        []hetznerSSHKeyModel `tfsdk:"ssh_keys"`
}

type hetznerSSHKeysDataSource struct{ client *coolify.Client }

func NewHetznerSSHKeysDataSource() datasource.DataSource { return &hetznerSSHKeysDataSource{} }

func (d *hetznerSSHKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_hetzner_ssh_keys"
}

func (d *hetznerSSHKeysDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists SSH keys uploaded to your Hetzner Cloud account. Use the key IDs when creating Hetzner servers via the Coolify API.",
		Attributes: map[string]schema.Attribute{
			"cloud_token_uuid": schema.StringAttribute{
				Required:    true,
				Description: "UUID of the Hetzner cloud token stored in Coolify.",
			},
			"ssh_keys": schema.ListNestedAttribute{
				Computed:    true,
				Description: "SSH keys available in the Hetzner account.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.Int64Attribute{Computed: true, Description: "Hetzner SSH key ID."},
						"name":        schema.StringAttribute{Computed: true, Description: "Key name."},
						"fingerprint": schema.StringAttribute{Computed: true, Description: "MD5 fingerprint."},
					},
				},
			},
		},
	}
}

func (d *hetznerSSHKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if client, ok := req.ProviderData.(*coolify.Client); ok {
		d.client = client
	}
}

func (d *hetznerSSHKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config hetznerSSHKeysDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var out any
	if err := d.client.Get(ctx, "/api/v1/hetzner/ssh-keys?cloud_token_uuid="+config.CloudTokenUUID.ValueString(), &out); err != nil {
		resp.Diagnostics.AddError("Unable to read Hetzner SSH keys", err.Error())
		return
	}
	config.SSHKeys = []hetznerSSHKeyModel{}
	for _, item := range dataList(out) {
		config.SSHKeys = append(config.SSHKeys, hetznerSSHKeyModel{
			ID:          types.Int64Value(int64FromAny(item["id"])),
			Name:        types.StringValue(stringFromAny(item["name"])),
			Fingerprint: types.StringValue(stringFromAny(item["fingerprint"])),
		})
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
