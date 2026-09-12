package provider

import (
	"encoding/json"
	"fmt"
)

func extractID(v json.RawMessage) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(v, &obj); err != nil {
		return ""
	}
	if raw, ok := obj["id"]; ok {
		return fmt.Sprintf("%v", raw)
	}
	if raw, ok := obj["_id"]; ok {
		return fmt.Sprintf("%v", raw)
	}
	return ""
}
