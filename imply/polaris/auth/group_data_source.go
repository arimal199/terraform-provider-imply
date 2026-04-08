package auth

import (
	"context"
	"fmt"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

func NewGroupDataSource() datasource.DataSource { return &groupDataSource{} }

type groupDataSource struct{ client *client.Client }

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := groupAttributes()
	attrs["id"] = schema.StringAttribute{Required: true}
	resp.Schema = schema.Schema{Attributes: attrs}
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.Get(fmt.Sprintf("/groups/%s", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Imply Group", err.Error())
		return
	}

	decoded := decodeGroup(group)
	decoded.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &decoded)...)
}

func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
