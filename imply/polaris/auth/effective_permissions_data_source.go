package auth

import (
	"context"
	"fmt"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &effectivePermissionsDataSource{}
	_ datasource.DataSourceWithConfigure = &effectivePermissionsDataSource{}
)

func NewEffectivePermissionsDataSource() datasource.DataSource {
	return &effectivePermissionsDataSource{}
}

type effectivePermissionsDataSource struct{ client *client.Client }

type effectivePermissionsModel struct {
	UserID types.String      `tfsdk:"user_id"`
	Items  []PermissionModel `tfsdk:"items"`
}

func (d *effectivePermissionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_effective_permissions"
}

func (d *effectivePermissionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"user_id": schema.StringAttribute{Required: true},
		"items":   schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: permissionAttributes()}},
	}}
}

func (d *effectivePermissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state effectivePermissionsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.client.Get(fmt.Sprintf("/users/%s/effectivepermissions", state.UserID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Imply Effective Permissions", err.Error())
		return
	}

	rawPermissions, ok := response["values"].([]any)
	if !ok {
		resp.Diagnostics.AddError("Invalid Response Format", fmt.Sprintf("Expected []any in values field, got: %T", response["values"]))
		return
	}

	state.Items = make([]PermissionModel, 0, len(rawPermissions))
	for _, p := range rawPermissions {
		permMap, ok := p.(map[string]any)
		if !ok {
			continue
		}
		state.Items = append(state.Items, decodePermission(permMap))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *effectivePermissionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}
	d.client = c
}
