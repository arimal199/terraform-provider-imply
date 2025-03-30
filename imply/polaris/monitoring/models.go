package monitoring

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AlertModel struct {
	ID              types.String  `tfsdk:"id"`
	Name            types.String  `tfsdk:"name"`
	Description     types.String  `tfsdk:"description"`
	Query           types.String  `tfsdk:"query"`
	Condition       types.String  `tfsdk:"condition"`
	Threshold       types.Float64 `tfsdk:"threshold"`
	Enabled         types.Bool    `tfsdk:"enabled"`
	CreatedOn       types.String  `tfsdk:"created_on"`
	LastTriggeredOn types.String  `tfsdk:"last_triggered_on"`
}
