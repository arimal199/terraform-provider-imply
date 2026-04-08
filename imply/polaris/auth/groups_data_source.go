// Copyright IBM Corp. 2026

package auth

import (
	"context"
	"fmt"

	"github.com/arimal199/terraform-provider-imply/imply/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &groupsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupsDataSource{}
)

func NewGroupsDataSource() datasource.DataSource { return &groupsDataSource{} }

type groupsDataSource struct{ client *client.Client }

func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{Attributes: map[string]schema.Attribute{
		"items": schema.ListNestedAttribute{Computed: true, NestedObject: schema.NestedAttributeObject{Attributes: groupAttributes()}},
	}}
}

func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GroupsModel
	response, err := d.client.Get("/groups")
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Imply Groups", err.Error())
		return
	}

	rawGroups, ok := response["values"].([]any)
	if !ok {
		resp.Diagnostics.AddError("Invalid Response Format", fmt.Sprintf("Expected []any in values field, got: %T", response["values"]))
		return
	}

	state.Items = make([]GroupModel, 0, len(rawGroups))
	for _, g := range rawGroups {
		groupMap, ok := g.(map[string]any)
		if !ok {
			continue
		}
		state.Items = append(state.Items, decodeGroup(groupMap))
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
