package auth

import (
	"context"
	"fmt"

	"github.com/arimal199/terraform-provider-imply/imply/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

type usersDataSource struct {
	client *client.Client
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						"username": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"first_name": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"last_name": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
						"email_verified": schema.BoolAttribute{
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
						"groups": schema.ListNestedAttribute{
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
						"identities": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Optional:    true,
						},
						"actions": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
							Optional:    true,
						},
						"created_on": schema.StringAttribute{
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
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UsersModel

	response, err := d.client.Get("/users")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Imply Users",
			err.Error(),
		)
		return
	}

	// Try to get users from the response
	users, ok := response["values"].([]any)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Response Format",
			fmt.Sprintf("Expected []any in values field, got: %T", response["values"]),
		)
		return
	}

	// Map response body to model
	for _, rawUser := range users {
		user, ok := rawUser.(map[string]interface{})
		if !ok {
			resp.Diagnostics.AddError(
				"Invalid User Data",
				fmt.Sprintf("Expected map[string]interface{}, got: %T", rawUser),
			)
			return
		}

		userState := UserModel{
			ID:        types.StringValue(fmt.Sprintf("%v", user["id"])),
			Username:  types.StringValue(fmt.Sprintf("%v", user["username"])),
			Email:     types.StringValue(fmt.Sprintf("%v", user["email"])),
			FirstName: types.StringValue(fmt.Sprintf("%v", user["firstName"])),
			LastName:  types.StringValue(fmt.Sprintf("%v", user["lastName"])),
		}

		// Handle boolean fields
		if enabled, ok := user["enabled"].(bool); ok {
			userState.Enabled = types.BoolValue(enabled)
		}
		if emailVerified, ok := user["emailVerified"].(bool); ok {
			userState.EmailVerified = types.BoolValue(emailVerified)
		}

		// Handle permissions
		if perms, ok := user["permissions"].([]any); ok && len(perms) > 0 {
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
				if resources, ok := perm["resources"].([]any); ok && len(resources) > 0 {
					for _, r := range resources {
						permModel.Resources = append(permModel.Resources, types.StringValue(fmt.Sprintf("%v", r)))
					}
					userState.Permissions = append(userState.Permissions, permModel)
				}
			}
		}

		// Handle groups
		if groups, ok := user["groups"].([]any); ok && len(groups) > 0 {
			for _, g := range groups {
				group, ok := g.(map[string]interface{})
				if !ok {
					continue
				}
				groupModel := GroupModel{
					ID:   types.StringValue(fmt.Sprintf("%v", group["id"])),
					Name: types.StringValue(fmt.Sprintf("%v", group["name"])),
				}

				if readOnly, ok := group["readOnly"].(bool); ok {
					groupModel.ReadOnly = types.BoolValue(readOnly)
				}
				if userCount, ok := group["userCount"].(float64); ok {
					groupModel.UserCount = types.Int64Value(int64(userCount))
				}

				// Handle group permissions
				if perms, ok := group["permissions"].([]any); ok && len(perms) > 0 {
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
						if resources, ok := perm["resources"].([]any); ok && len(resources) > 0 {
							for _, r := range resources {
								permModel.Resources = append(permModel.Resources, types.StringValue(fmt.Sprintf("%v", r)))
							}
							groupModel.Permissions = append(groupModel.Permissions, permModel)
						}
					}
				}

				userState.Groups = append(userState.Groups, groupModel)
			}
		}

		// Handle arrays
		if identities, ok := user["identities"].([]any); ok && len(identities) > 0 {
			for _, identity := range identities {
				userState.Identities = append(userState.Identities, types.StringValue(fmt.Sprintf("%v", identity)))
			}
		}

		if actions, ok := user["actions"].([]any); ok && len(actions) > 0 {
			for _, action := range actions {
				userState.Actions = append(userState.Actions, types.StringValue(fmt.Sprintf("%v", action)))
			}
		}

		// Handle timestamps
		if createdOn, ok := user["createdOn"].(string); ok {
			userState.CreatedOn = types.StringValue(createdOn)
		}

		state.Items = append(state.Items, userState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
