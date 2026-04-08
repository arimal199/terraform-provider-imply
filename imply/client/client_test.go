package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNormalizeHost(t *testing.T) {
	h := "https://x.app.imply.io/"
	if got := normalizeHost(h); got != "https://x.api.imply.io/v1" {
		t.Fatalf("unexpected normalized host: %s", got)
	}
}

func TestDoRequestHeadersAndJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Basic token" {
			t.Fatalf("missing auth header: %s", got)
		}
		if r.Method != http.MethodGet || r.URL.Path != "/v1/users" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"values":[]}`))
	}))
	defer ts.Close()

	host := ts.URL
	key := "token"
	c, err := NewClient(&host, &key)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Get("/users")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resp["values"]; !ok {
		t.Fatal("expected values")
	}
}

func TestNoContentDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	host := ts.URL
	key := "token"
	c, _ := NewClient(&host, &key)
	if err := c.Delete("/users/1"); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteWithBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != `{"id":"u1"}` {
			t.Fatalf("unexpected request body: %s", string(body))
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()
	host := ts.URL
	key := "token"
	c, _ := NewClient(&host, &key)
	if _, err := c.DeleteWithBody("/groups/g1/members", map[string]string{"id": "u1"}); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer ts.Close()
	host := ts.URL
	key := "token"
	c, _ := NewClient(&host, &key)
	if _, err := c.Get("/bad"); err == nil {
		t.Fatal("expected error")
	}
}
