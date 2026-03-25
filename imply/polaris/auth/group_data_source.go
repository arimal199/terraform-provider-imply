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
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	client *client.Client
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
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
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group by ID
	group, err := d.client.Get(fmt.Sprintf("/groups/%s", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Imply Group",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(fmt.Sprintf("%v", group["id"]))
	state.Name = types.StringValue(fmt.Sprintf("%v", group["name"]))

	// Handle read_only
	if readOnly, ok := group["readOnly"].(bool); ok {
		state.ReadOnly = types.BoolValue(readOnly)
	}

	// Handle user_count
	if userCount, ok := group["userCount"].(float64); ok {
		state.UserCount = types.Int64Value(int64(userCount))
	}

	// Handle permissions
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
				state.Permissions = append(state.Permissions, permModel)
			}
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
