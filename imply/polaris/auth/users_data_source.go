// Copyright IBM Corp. 2026

package auth

import (
	"context"
	"fmt"
	"net/url"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

func NewUsersDataSource() datasource.DataSource { return &usersDataSource{} }

type usersDataSource struct{ client *client.Client }

type usersDataSourceModel struct {
	Top    types.Int64  `tfsdk:"top"`
	Skip   types.Int64  `tfsdk:"skip"`
	Search types.String `tfsdk:"search"`
	Items  []UserModel  `tfsdk:"items"`
}

func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"top":    schema.Int64Attribute{Optional: true},
		"skip":   schema.Int64Attribute{Optional: true},
		"search": schema.StringAttribute{Optional: true},
		"items":  schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: userAttributes()}},
	}}
}

func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	q := url.Values{}
	if !state.Top.IsNull() {
		q.Set("top", fmt.Sprintf("%d", state.Top.ValueInt64()))
	}
	if !state.Skip.IsNull() {
		q.Set("skip", fmt.Sprintf("%d", state.Skip.ValueInt64()))
	}
	if !state.Search.IsNull() && state.Search.ValueString() != "" {
		q.Set("search", state.Search.ValueString())
	}
	path := "/users"
	if len(q) > 0 {
		path = path + "?" + q.Encode()
	}

	response, err := d.client.Get(path)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Imply Users", err.Error())
		return
	}

	rawUsers, ok := response["values"].([]any)
	if !ok {
		resp.Diagnostics.AddError("Invalid Response Format", fmt.Sprintf("Expected []any in values field, got: %T", response["values"]))
		return
	}

	state.Items = make([]UserModel, 0, len(rawUsers))
	for _, rawUser := range rawUsers {
		userMap, ok := rawUser.(map[string]any)
		if !ok {
			continue
		}
		state.Items = append(state.Items, decodeUser(userMap))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
