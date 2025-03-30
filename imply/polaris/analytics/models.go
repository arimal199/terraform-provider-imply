package analytics

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataCubeModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Table       types.String `tfsdk:"table"`
	Config      types.String `tfsdk:"config"`
	CreatedOn   types.String `tfsdk:"created_on"`
	LastUsedOn  types.String `tfsdk:"last_used_on"`
}

type DashboardModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Layout      types.String `tfsdk:"layout"`
	CreatedOn   types.String `tfsdk:"created_on"`
	LastUsedOn  types.String `tfsdk:"last_used_on"`
}
