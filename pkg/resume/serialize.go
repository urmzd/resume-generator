package resume

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// SerializeResume marshals a Resume to bytes in the given format.
// Returns the serialized bytes, the canonical format name, and any error.
func SerializeResume(r *Resume, format string) ([]byte, string, error) {
	switch strings.ToLower(format) {
	case "yaml", "yml":
		data, err := yaml.Marshal(r)
		return data, "yaml", err
	case "json":
		data, err := json.MarshalIndent(r, "", "  ")
		return data, "json", err
	case "toml":
		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(r)
		return buf.Bytes(), "toml", err
	case "md", "markdown":
		// Markdown input cannot be losslessly serialized; fall back to YAML
		data, err := yaml.Marshal(r)
		return data, "yaml", err
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}
}
