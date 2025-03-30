package auth

import (
	"context"
	"fmt"

	"github.com/arimal199/terraform-provider-imply/imply/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &permissionsDataSource{}
	_ datasource.DataSourceWithConfigure = &permissionsDataSource{}
)

type permissionsDataSourceModel struct {
	Items []PermissionModel `tfsdk:"items"`
}

func NewPermissionsDataSource() datasource.DataSource {
	return &permissionsDataSource{}
}

type permissionsDataSource struct {
	client *client.Client
}

func (d *permissionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permissions"
}

func (d *permissionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"resources": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *permissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state permissionsDataSourceModel

	response, err := d.client.Get("/permissions")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Imply Permissions",
			err.Error(),
		)
		return
	}

	// Try to get permissions from the response
	permissions, ok := response["values"].([]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Response Format",
			fmt.Sprintf("Expected []interface{} in values field, got: %T", response["values"]),
		)
		return
	}

	// Map response body to model
	for _, rawPermission := range permissions {
		permission, ok := rawPermission.(map[string]interface{})
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid Permission Data",
				fmt.Sprintf("Expected map[string]interface{}, got: %T", rawPermission),
			)
			return
		}

		permissionState := PermissionModel{
			ID:   types.StringValue(fmt.Sprintf("%v", permission["id"])),
			Name: types.StringValue(fmt.Sprintf("%v", permission["name"])),
		}

		// Handle resources
		if resources, ok := permission["resources"].([]interface{}); ok && len(resources) > 0 {
			for _, r := range resources {
				permissionState.Resources = append(permissionState.Resources, types.StringValue(fmt.Sprintf("%v", r)))
			}
		}

		state.Items = append(state.Items, permissionState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *permissionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}
