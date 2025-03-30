package auth

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PermissionModel struct {
	ID        types.String   `tfsdk:"id"`
	Name      types.String   `tfsdk:"name"`
	Resources []types.String `tfsdk:"resources"`
}

type GroupModel struct {
	ID          types.String      `tfsdk:"id"`
	Name        types.String      `tfsdk:"name"`
	ReadOnly    types.Bool        `tfsdk:"read_only"`
	Permissions []PermissionModel `tfsdk:"permissions"`
	UserCount   types.Int64       `tfsdk:"user_count"`
}

type UserModel struct {
	ID            types.String      `tfsdk:"id"`
	Username      types.String      `tfsdk:"username"`
	Email         types.String      `tfsdk:"email"`
	FirstName     types.String      `tfsdk:"first_name"`
	LastName      types.String      `tfsdk:"last_name"`
	Enabled       types.Bool        `tfsdk:"enabled"`
	EmailVerified types.Bool        `tfsdk:"email_verified"`
	Permissions   []PermissionModel `tfsdk:"permissions"`
	Groups        []GroupModel      `tfsdk:"groups"`
	Identities    []types.String    `tfsdk:"identities"`
	Actions       []types.String    `tfsdk:"actions"`
	CreatedOn     types.String      `tfsdk:"created_on"`
}
