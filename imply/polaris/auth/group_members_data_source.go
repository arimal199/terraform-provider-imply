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
	_ datasource.DataSource              = &groupMembersDataSource{}
	_ datasource.DataSourceWithConfigure = &groupMembersDataSource{}
)

func NewGroupMembersDataSource() datasource.DataSource { return &groupMembersDataSource{} }

type groupMembersDataSource struct{ client *client.Client }

type groupMembersModel struct {
	GroupID types.String `tfsdk:"group_id"`
	Top     types.Int64  `tfsdk:"top"`
	Skip    types.Int64  `tfsdk:"skip"`
	Search  types.String `tfsdk:"search"`
	Items   []UserModel  `tfsdk:"items"`
}

func (d *groupMembersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_members"
}

func (d *groupMembersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"group_id": schema.StringAttribute{Required: true},
		"top":      schema.Int64Attribute{Optional: true},
		"skip":     schema.Int64Attribute{Optional: true},
		"search":   schema.StringAttribute{Optional: true},
		"items":    schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: userAttributes()}},
	}}
}

func (d *groupMembersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupMembersModel
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

	path := fmt.Sprintf("/groups/%s/members", state.GroupID.ValueString())
	if len(q) > 0 {
		path = path + "?" + q.Encode()
	}

	response, err := d.client.Get(path)
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Imply Group Members", err.Error())
		return
	}

	rawUsers, ok := response["values"].([]any)
	if !ok {
		resp.Diagnostics.AddError("Invalid Response Format", fmt.Sprintf("Expected []any in values field, got: %T", response["values"]))
		return
	}

	state.Items = make([]UserModel, 0, len(rawUsers))
	for _, u := range rawUsers {
		userMap, ok := u.(map[string]any)
		if !ok {
			continue
		}
		state.Items = append(state.Items, decodeUser(userMap))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *groupMembersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
