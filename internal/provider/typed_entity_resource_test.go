package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestHelpcenterResourceSetsCnameDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", request.Method, http.MethodGet)
		}
		if request.URL.Path != "/helpcenters/123" {
			t.Errorf("request path = %q, want %q", request.URL.Path, "/helpcenters/123")
		}
		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`{"id":123,"name":"Help","language":"en-US","public_url":"https://help-e6ec79.herodesk-help.io"}`))
	}))
	defer server.Close()

	resource := resourceHelpcenter()
	data := schema.TestResourceDataRaw(t, resource.Schema, map[string]interface{}{
		"name":     "Help",
		"language": "en-US",
	})
	data.SetId("123")
	client := &Client{baseURL: server.URL, token: "test-token", httpClient: server.Client()}

	if diags := resource.ReadContext(context.Background(), data, client); diags.HasError() {
		t.Fatalf("read returned diagnostics: %v", diags)
	}
	if cnameDomain := data.Get("cname_domain"); cnameDomain != "help-e6ec79.herodesk-help.io" {
		t.Errorf("resource cname_domain = %q, want %q", cnameDomain, "help-e6ec79.herodesk-help.io")
	}
}
