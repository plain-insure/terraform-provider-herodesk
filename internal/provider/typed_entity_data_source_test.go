package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestTypedDataSourceLookupFieldDoesNotConflictWithStateID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			t.Errorf("request method = %q, want %q", request.Method, http.MethodGet)
		}
		if request.URL.Path != "/helpcenters/123" {
			t.Errorf("request path = %q, want %q", request.URL.Path, "/helpcenters/123")
		}
		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`{"id":123,"name":"Help","language":"en-US"}`))
	}))
	defer server.Close()

	dataSource := dataSourceHelpcenter()
	data := schema.TestResourceDataRaw(t, dataSource.Schema, map[string]interface{}{"helpcenter_id": 123})
	client := &Client{baseURL: server.URL, token: "test-token", httpClient: server.Client()}

	if diags := dataSource.ReadContext(context.Background(), data, client); diags.HasError() {
		t.Fatalf("read returned diagnostics: %v", diags)
	}
	if data.Id() != "helpcenters-123" {
		t.Errorf("data source state ID = %q, want %q", data.Id(), "helpcenters-123")
	}
	_ = data.State()
}
