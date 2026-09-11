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

func TestFlattenItems(t *testing.T) {
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
			items := flattenItems([]byte(tt.json))
			if len(items) != tt.want {
				t.Fatalf("flattenItems() len = %d, want %d", len(items), tt.want)
			}
		})
	}
}
