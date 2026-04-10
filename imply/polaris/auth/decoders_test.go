// Copyright IBM Corp. 2026

package auth

import "testing"

func TestDecodeUser(t *testing.T) {
	in := map[string]any{
		"id": "u1", "username": "alice", "email": "a@example.com", "enabled": true,
		"permissions": []any{map[string]any{"id": "p1", "name": "read", "resources": []any{"*"}}},
	}
	u := decodeUser(in)
	if u.ID.ValueString() != "u1" || u.Username.ValueString() != "alice" {
		t.Fatalf("unexpected user decode: %#v", u)
	}
	if len(u.Permissions) != 1 {
		t.Fatalf("expected permission")
	}
}

func TestDecodeUserEmpty(t *testing.T) {
	u := decodeUser(map[string]any{})
	if !u.ID.IsNull() {
		t.Fatalf("expected null id")
	}
}
