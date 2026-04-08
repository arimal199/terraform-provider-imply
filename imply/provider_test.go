package imply

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResolveProviderConfig(t *testing.T) {
	t.Setenv("IMPLY_HOST", "env-host")
	t.Setenv("IMPLY_API_KEY", "env-key")

	host, key := resolveProviderConfig(implyProviderModel{})
	if host != "env-host" || key != "env-key" {
		t.Fatalf("unexpected env fallback: %s %s", host, key)
	}

	host, key = resolveProviderConfig(implyProviderModel{
		Host:   types.StringValue("cfg-host"),
		ApiKey: types.StringValue("cfg-key"),
	})
	if host != "cfg-host" || key != "cfg-key" {
		t.Fatalf("unexpected config override: %s %s", host, key)
	}
}

func TestResolveProviderConfigEmptyEnv(t *testing.T) {
	_ = os.Unsetenv("IMPLY_HOST")
	_ = os.Unsetenv("IMPLY_API_KEY")
	h, k := resolveProviderConfig(implyProviderModel{})
	if h != "" || k != "" {
		t.Fatalf("expected empty values")
	}
}
