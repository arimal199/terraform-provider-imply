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
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource { return &userDataSource{} }

type userDataSource struct{ client *client.Client }

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := userAttributes()
	attrs["id"] = schema.StringAttribute{Required: true}
	resp.Schema = schema.Schema{Attributes: attrs}
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.Get(fmt.Sprintf("/users/%s", state.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Imply User", err.Error())
		return
	}

	decoded := decodeUser(user)
	decoded.ID = types.StringValue(state.ID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &decoded)...)
}

func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
