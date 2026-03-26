// Copyright (c) HashiCorp, Inc.

package auth

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "status: 404")
}

func stringValue(data map[string]any, key string) types.String {
	value, ok := data[key]
	if !ok || value == nil {
		return types.StringNull()
	}

	str := fmt.Sprintf("%v", value)
	if str == "" || str == "<nil>" {
		return types.StringNull()
	}

	return types.StringValue(str)
}

func boolValue(data map[string]any, key string) types.Bool {
	value, ok := data[key].(bool)
	if !ok {
		return types.BoolNull()
	}

	return types.BoolValue(value)
}

func int64Value(data map[string]any, key string) types.Int64 {
	value, ok := data[key].(float64)
	if !ok {
		return types.Int64Null()
	}

	return types.Int64Value(int64(value))
}

func permissionModels(raw any) []PermissionModel {
	values, ok := raw.([]any)
	if !ok || len(values) == 0 {
		return []PermissionModel{}
	}

	items := make([]PermissionModel, 0, len(values))
	for _, value := range values {
		permission, ok := value.(map[string]any)
		if !ok {
			continue
		}

		item := PermissionModel{
			ID:   stringValue(permission, "id"),
			Name: stringValue(permission, "name"),
		}

		if resources, ok := permission["resources"].([]any); ok {
			item.Resources = make([]types.String, 0, len(resources))
			for _, resource := range resources {
				item.Resources = append(item.Resources, types.StringValue(fmt.Sprintf("%v", resource)))
			}
		}

		items = append(items, item)
	}

	return items
}

func groupModels(raw any) []GroupModel {
	values, ok := raw.([]any)
	if !ok || len(values) == 0 {
		return []GroupModel{}
	}

	items := make([]GroupModel, 0, len(values))
	for _, value := range values {
		group, ok := value.(map[string]any)
		if !ok {
			continue
		}

		items = append(items, GroupModel{
			ID:          stringValue(group, "id"),
			Name:        stringValue(group, "name"),
			ReadOnly:    boolValue(group, "readOnly"),
			Permissions: permissionModels(group["permissions"]),
			UserCount:   int64Value(group, "userCount"),
		})
	}

	return items
}

func stringModels(raw any, key string) []types.String {
	values, ok := raw.([]any)
	if !ok || len(values) == 0 {
		return []types.String{}
	}

	items := make([]types.String, 0, len(values))
	for _, value := range values {
		switch typed := value.(type) {
		case string:
			items = append(items, types.StringValue(typed))
		case map[string]any:
			if providerID, ok := typed[key].(string); ok && providerID != "" {
				items = append(items, types.StringValue(providerID))
			}
		default:
			items = append(items, types.StringValue(fmt.Sprintf("%v", value)))
		}
	}

	return items
}

func groupMemberExists(raw any, userID string) bool {
	values, ok := raw.([]any)
	if !ok {
		return false
	}

	for _, value := range values {
		member, ok := value.(map[string]any)
		if !ok {
			continue
		}

		if fmt.Sprintf("%v", member["id"]) == userID {
			return true
		}
	}

	return false
}
