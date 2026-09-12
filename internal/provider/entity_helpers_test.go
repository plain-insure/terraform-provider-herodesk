package provider

import "testing"

func TestExtractID(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{name: "id", json: `{"id":"abc"}`, want: "abc"},
		{name: "underscore id", json: `{"_id":"xyz"}`, want: "xyz"},
		{name: "missing", json: `{"name":"n"}`, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractID([]byte(tt.json)); got != tt.want {
				t.Fatalf("extractID() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCnameDomain(t *testing.T) {
	tests := []struct {
		name      string
		publicURL string
		want      string
	}{
		{name: "host", publicURL: "https://help-e6ec79.herodesk-help.io", want: "help-e6ec79.herodesk-help.io"},
		{name: "host with port and path", publicURL: "https://help.example.com:8443/articles/1", want: "help.example.com"},
		{name: "invalid URL", publicURL: "://invalid", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cnameDomain(map[string]interface{}{"public_url": tt.publicURL}); got != tt.want {
				t.Fatalf("cnameDomain() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResponseObjects(t *testing.T) {
	tests := []struct {
		name string
		json string
		want int
	}{
		{name: "array", json: `[{"id":1},{"id":2}]`, want: 2},
		{name: "items key", json: `{"items":[{"id":1}]}`, want: 1},
		{name: "single object", json: `{"id":1}`, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := responseObjects([]byte(tt.json))
			if err != nil {
				t.Fatalf("responseObjects() returned an error: %v", err)
			}
			if len(items) != tt.want {
				t.Fatalf("responseObjects() len = %d, want %d", len(items), tt.want)
			}
		})
	}
}
