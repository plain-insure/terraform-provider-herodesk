package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientUpdateUsesPatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPatch {
			t.Errorf("request method = %q, want %q", request.Method, http.MethodPatch)
		}
		if request.URL.Path != "/tags/123" {
			t.Errorf("request path = %q, want %q", request.URL.Path, "/tags/123")
		}
		responseWriter.Header().Set("Content-Type", "application/json")
		_, _ = responseWriter.Write([]byte(`{"id":123}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: server.Client(),
	}

	if _, err := client.update(context.Background(), "tags", "123", []byte(`{"name":"Priority"}`)); err != nil {
		t.Fatalf("update() returned an error: %v", err)
	}
}