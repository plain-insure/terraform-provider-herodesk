package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestHelpcenterCustomDomainResourcePatchesOnlyCustomDomain(t *testing.T) {
	var customDomains []string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/helpcenters/123" {
			t.Errorf("request path = %q, want %q", request.URL.Path, "/helpcenters/123")
		}
		switch request.Method {
		case http.MethodPatch:
			var payload map[string]string
			if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
				t.Fatalf("decode patch payload: %v", err)
			}
			if len(payload) != 1 {
				t.Errorf("patch payload = %#v, want only custom_domain", payload)
			}
			customDomains = append(customDomains, payload["custom_domain"])
			responseWriter.Header().Set("Content-Type", "application/json")
			_, _ = responseWriter.Write([]byte(`{"id":123}`))
		case http.MethodGet:
			responseWriter.Header().Set("Content-Type", "application/json")
			if len(customDomains) == 0 {
				t.Fatal("GET called before PATCH")
			}
			_, _ = responseWriter.Write([]byte(`{"id":123,"custom_domain":"docs.example.com"}`))
		default:
			t.Errorf("request method = %q, want PATCH or GET", request.Method)
		}
	}))
	defer server.Close()

	resource := resourceHelpcenterCustomDomain()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"helpcenter_id": 123,
		"custom_domain": "docs.example.com",
	})
	client := &Client{baseURL: server.URL, token: "test-token", httpClient: server.Client()}

	if diags := resource.CreateContext(context.Background(), data, client); diags.HasError() {
		t.Fatalf("create returned diagnostics: %v", diags)
	}
	if data.Id() != "123" {
		t.Errorf("resource ID = %q, want %q", data.Id(), "123")
	}
	if diags := resource.DeleteContext(context.Background(), data, client); diags.HasError() {
		t.Fatalf("delete returned diagnostics: %v", diags)
	}

	if len(customDomains) != 2 {
		t.Fatalf("PATCH count = %d, want 2", len(customDomains))
	}
	if customDomains[0] != "docs.example.com" {
		t.Errorf("create custom domain = %q, want %q", customDomains[0], "docs.example.com")
	}
	if customDomains[1] != "" {
		t.Errorf("delete custom domain = %q, want empty string", customDomains[1])
	}
}
