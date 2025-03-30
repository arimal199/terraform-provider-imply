package data

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ColumnModel struct {
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Nullable    types.Bool   `tfsdk:"nullable"`
	Description types.String `tfsdk:"description"`
}

type TableModel struct {
	ID             types.String  `tfsdk:"id"`
	Name           types.String  `tfsdk:"name"`
	Type           types.String  `tfsdk:"type"`
	SchemaMode     types.String  `tfsdk:"schema_mode"`
	Partitioning   types.String  `tfsdk:"partitioning"`
	Rollup         types.Bool    `tfsdk:"rollup"`
	CreatedOn      types.String  `tfsdk:"created_on"`
	LastModifiedOn types.String  `tfsdk:"last_modified_on"`
	RowCount       types.Int64   `tfsdk:"row_count"`
	SizeBytes      types.Int64   `tfsdk:"size_bytes"`
	Columns        []ColumnModel `tfsdk:"columns"`
}

type ConnectionModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	Config      types.String `tfsdk:"config"`
	CreatedOn   types.String `tfsdk:"created_on"`
	LastUsedOn  types.String `tfsdk:"last_used_on"`
}

type JobModel struct {
	ID            types.String  `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	Type          types.String  `tfsdk:"type"`
	Status        types.String  `tfsdk:"status"`
	Source        types.String  `tfsdk:"source"`
	Destination   types.String  `tfsdk:"destination"`
	Config        types.String  `tfsdk:"config"`
	CreatedOn     types.String  `tfsdk:"created_on"`
	StartedOn     types.String  `tfsdk:"started_on"`
	CompletedOn   types.String  `tfsdk:"completed_on"`
	Error         types.String  `tfsdk:"error"`
	Progress      types.Float64 `tfsdk:"progress"`
	RowsProcessed types.Int64   `tfsdk:"rows_processed"`
}
