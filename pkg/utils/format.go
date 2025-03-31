package utils

import (
	"encoding/json"
	"fmt"
)

// FormatValue converts a value to a readable string representation
func FormatValue(v any) string {
	if v == nil {
		return "<nil>"
	}

	switch value := v.(type) {
	case string:
		return Truncate(value, 40)
	case []any:
		b, err := json.Marshal(value)
		if err != nil {
			return fmt.Sprintf("%v", value)
		}
		return Truncate(string(b), 40)
	case map[string]any:
		b, err := json.Marshal(value)
		if err != nil {
			return fmt.Sprintf("%v", value)
		}
		return Truncate(string(b), 40)
	default:
		return fmt.Sprintf("%v", value)
	}
}
