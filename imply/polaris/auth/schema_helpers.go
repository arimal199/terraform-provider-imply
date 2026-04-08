package auth

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func permissionAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":        schema.StringAttribute{Computed: true, Optional: true},
		"name":      schema.StringAttribute{Computed: true, Optional: true},
		"resources": schema.ListAttribute{Computed: true, Optional: true, ElementType: types.StringType},
	}
}

func groupAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":          schema.StringAttribute{Computed: true, Optional: true},
		"name":        schema.StringAttribute{Computed: true, Optional: true},
		"read_only":   schema.BoolAttribute{Computed: true, Optional: true},
		"permissions": schema.ListNestedAttribute{Computed: true, Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: permissionAttributes()}},
		"user_count":  schema.Int64Attribute{Computed: true, Optional: true},
	}
}

func userAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":             schema.StringAttribute{Computed: true, Optional: true},
		"username":       schema.StringAttribute{Computed: true, Optional: true},
		"email":          schema.StringAttribute{Computed: true, Optional: true},
		"first_name":     schema.StringAttribute{Computed: true, Optional: true},
		"last_name":      schema.StringAttribute{Computed: true, Optional: true},
		"enabled":        schema.BoolAttribute{Computed: true, Optional: true},
		"email_verified": schema.BoolAttribute{Computed: true, Optional: true},
		"permissions":    schema.ListNestedAttribute{Computed: true, Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: permissionAttributes()}},
		"groups":         schema.ListNestedAttribute{Computed: true, Optional: true, NestedObject: schema.NestedAttributeObject{Attributes: groupAttributes()}},
		"identities":     schema.ListAttribute{Computed: true, Optional: true, ElementType: types.StringType},
		"actions":        schema.ListAttribute{Computed: true, Optional: true, ElementType: types.StringType},
		"created_on":     schema.StringAttribute{Computed: true, Optional: true},
	}
}
