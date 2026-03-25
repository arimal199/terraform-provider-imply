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
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *client.Client
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"username": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"first_name": schema.StringAttribute{
				Computed: true,
			},
			"last_name": schema.StringAttribute{
				Computed: true,
			},
			"enabled": schema.BoolAttribute{
				Computed: true,
			},
			"email_verified": schema.BoolAttribute{
				Computed: true,
			},
			"permissions": schema.ListNestedAttribute{
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
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"read_only": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
						"permissions": schema.ListNestedAttribute{
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
			},
			"actions": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_on": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserModel
	// var config UserModel

	// Read Terraform configuration data into the model
	// resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get user by ID
	// user, err := d.client.Get(fmt.Sprintf("/users/%s", config.ID.ValueString()))
	user, err := d.client.Get(fmt.Sprintf("/users/%s", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Imply User",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%v", user["id"]))
	state.Username = types.StringValue(fmt.Sprintf("%v", user["username"]))
	state.Email = types.StringValue(fmt.Sprintf("%v", user["email"]))

	if firstName, ok := user["firstName"].(string); ok {
		if firstName == "" {
			state.FirstName = types.StringNull()
		} else {
			state.FirstName = types.StringValue(fmt.Sprintf("%v", user["firstName"]))
		}
	}
	if lastName, ok := user["lastName"].(string); ok {
		if lastName == "" {
			state.LastName = types.StringNull()
		} else {
			state.LastName = types.StringValue(fmt.Sprintf("%v", user["lastName"]))
		}
	}

	// Handle boolean fields
	if enabled, ok := user["enabled"].(bool); ok {
		state.Enabled = types.BoolValue(enabled)
	}
	if emailVerified, ok := user["emailVerified"].(bool); ok {
		state.EmailVerified = types.BoolValue(emailVerified)
	}

	// Handle permissions
	if perms, ok := user["permissions"].([]any); ok {
		if len(perms) > 0 {
			for _, p := range perms {
				perm, ok := p.(map[string]any)
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
					state.Permissions = append(state.Permissions, permModel)
				}
			}
		} else {
			state.Permissions = []PermissionModel{}
		}
	}

	// Handle groups
	if groups, ok := user["groups"].([]any); ok {
		if len(groups) > 0 {
			for _, g := range groups {
				group, ok := g.(map[string]any)
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
						perm, ok := p.(map[string]any)
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

				state.Groups = append(state.Groups, groupModel)
			}
		} else {
			state.Groups = []GroupModel{}
		}
	}

	// Handle arrays
	if identities, ok := user["identities"].([]any); ok {
		if len(identities) > 0 {
			for _, identity := range identities {
				state.Identities = append(state.Identities, types.StringValue(fmt.Sprintf("%v", identity)))
			}
		} else {
			state.Identities = []types.String{}
		}
	}

	if actions, ok := user["actions"].([]any); ok && len(actions) > 0 {
		for _, action := range actions {
			state.Actions = append(state.Actions, types.StringValue(fmt.Sprintf("%v", action)))
		}
	}

	// Handle timestamps
	if createdOn, ok := user["createdOn"].(string); ok {
		state.CreatedOn = types.StringValue(createdOn)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
