// Copyright IBM Corp. 2026

package auth

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func decodePermission(input map[string]any) PermissionModel {
	perm := PermissionModel{
		ID:   stringFromAny(input["id"]),
		Name: stringFromAny(input["name"]),
	}

	if resources, ok := input["resources"].([]any); ok {
		perm.Resources = make([]types.String, 0, len(resources))
		for _, r := range resources {
			perm.Resources = append(perm.Resources, types.StringValue(fmt.Sprintf("%v", r)))
		}
	}

	return perm
}

func decodeGroup(input map[string]any) GroupModel {
	group := GroupModel{
		ID:   stringFromAny(input["id"]),
		Name: stringFromAny(input["name"]),
	}

	if readOnly, ok := input["readOnly"].(bool); ok {
		group.ReadOnly = types.BoolValue(readOnly)
	}
	if userCount, ok := input["userCount"].(float64); ok {
		group.UserCount = types.Int64Value(int64(userCount))
	}

	if perms, ok := input["permissions"].([]any); ok {
		group.Permissions = make([]PermissionModel, 0, len(perms))
		for _, p := range perms {
			permMap, ok := p.(map[string]any)
			if !ok {
				continue
			}
			group.Permissions = append(group.Permissions, decodePermission(permMap))
		}
	}

	return group
}

func decodeUser(input map[string]any) UserModel {
	user := UserModel{
		ID:       stringFromAny(input["id"]),
		Username: stringFromAny(input["username"]),
		Email:    stringFromAny(input["email"]),
	}

	if firstName, ok := input["firstName"].(string); ok && firstName != "" {
		user.FirstName = types.StringValue(firstName)
	} else {
		user.FirstName = types.StringNull()
	}

	if lastName, ok := input["lastName"].(string); ok && lastName != "" {
		user.LastName = types.StringValue(lastName)
	} else {
		user.LastName = types.StringNull()
	}

	if enabled, ok := input["enabled"].(bool); ok {
		user.Enabled = types.BoolValue(enabled)
	}
	if emailVerified, ok := input["emailVerified"].(bool); ok {
		user.EmailVerified = types.BoolValue(emailVerified)
	}

	if perms, ok := input["permissions"].([]any); ok {
		user.Permissions = make([]PermissionModel, 0, len(perms))
		for _, p := range perms {
			permMap, ok := p.(map[string]any)
			if !ok {
				continue
			}
			user.Permissions = append(user.Permissions, decodePermission(permMap))
		}
	}

	if groups, ok := input["groups"].([]any); ok {
		user.Groups = make([]GroupModel, 0, len(groups))
		for _, g := range groups {
			groupMap, ok := g.(map[string]any)
			if !ok {
				continue
			}
			user.Groups = append(user.Groups, decodeGroup(groupMap))
		}
	}

	if identities, ok := input["identities"].([]any); ok {
		user.Identities = make([]types.String, 0, len(identities))
		for _, identity := range identities {
			user.Identities = append(user.Identities, types.StringValue(fmt.Sprintf("%v", identity)))
		}
	}

	if actions, ok := input["actions"].([]any); ok {
		user.Actions = make([]types.String, 0, len(actions))
		for _, action := range actions {
			user.Actions = append(user.Actions, types.StringValue(fmt.Sprintf("%v", action)))
		}
	}

	if createdOn, ok := input["createdOn"].(string); ok {
		user.CreatedOn = types.StringValue(createdOn)
	}

	return user
}

func stringFromAny(v any) types.String {
	if v == nil {
		return types.StringNull()
	}
	s := fmt.Sprintf("%v", v)
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}
