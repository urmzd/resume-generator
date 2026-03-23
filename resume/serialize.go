package resume

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SerializeResume marshals a Resume to bytes in the given format.
// Returns the serialized bytes, the canonical format name, and any error.
func SerializeResume(r *Resume, format string) ([]byte, string, error) {
	switch strings.ToLower(format) {
	case "json":
		data, err := json.MarshalIndent(r, "", "  ")
		return data, "json", err
	case "md", "markdown":
		data, err := json.MarshalIndent(r, "", "  ")
		return data, "json", err
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}
}
