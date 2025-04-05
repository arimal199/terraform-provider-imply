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
	_ datasource.DataSource              = &groupsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupsDataSource{}
)


func NewGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

type groupsDataSource struct {
	client *client.Client
}

func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"read_only": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
						"permissions": schema.ListNestedAttribute{
							Computed: true,
							Optional: true,
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
						"user_count": schema.Int64Attribute{
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupsModel

	response, err := d.client.Get("/groups")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Imply Groups",
			err.Error(),
		)
		return
	}

	// Try to get groups from the response
	groups, ok := response["values"].([]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Response Format",
			fmt.Sprintf("Expected []interface{} in values field, got: %T", response["values"]),
		)
		return
	}

	// Map response body to model
	for _, rawGroup := range groups {
		group, ok := rawGroup.(map[string]interface{})
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid Group Data",
				fmt.Sprintf("Expected map[string]interface{}, got: %T", rawGroup),
			)
			return
		}

		groupState := GroupModel{
			ID:   types.StringValue(fmt.Sprintf("%v", group["id"])),
			Name: types.StringValue(fmt.Sprintf("%v", group["name"])),
		}

		// Handle read_only
		if readOnly, ok := group["readOnly"].(bool); ok {
			groupState.ReadOnly = types.BoolValue(readOnly)
		}

		// Handle user_count
		if userCount, ok := group["userCount"].(float64); ok {
			groupState.UserCount = types.Int64Value(int64(userCount))
		}

		// Handle permissions
		if perms, ok := group["permissions"].([]interface{}); ok && len(perms) > 0 {
			for _, p := range perms {
				perm, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				permModel := PermissionModel{
					ID:   types.StringValue(fmt.Sprintf("%v", perm["id"])),
					Name: types.StringValue(fmt.Sprintf("%v", perm["name"])),
				}

				// Handle resources
				if resources, ok := perm["resources"].([]interface{}); ok && len(resources) > 0 {
					for _, r := range resources {
						permModel.Resources = append(permModel.Resources, types.StringValue(fmt.Sprintf("%v", r)))
					}
					groupState.Permissions = append(groupState.Permissions, permModel)
				}
			}
		}

		state.Items = append(state.Items, groupState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
