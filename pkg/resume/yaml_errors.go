package resume

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type yamlPathPartKind int

const (
	yamlPathKey yamlPathPartKind = iota
	yamlPathIndex
)

type yamlPathPart struct {
	kind  yamlPathPartKind
	key   string
	index int
}

type yamlPathEntry struct {
	line  int
	col   int
	depth int
	path  string
}

var yamlLineRe = regexp.MustCompile(`line\s+(\d+)`)

// UnmarshalYAMLWithContext parses YAML and augments errors with line and path details.
func UnmarshalYAMLWithContext(data []byte, out any) error {
	if err := yaml.Unmarshal(data, out); err != nil {
		return formatYAMLError(data, err)
	}
	return nil
}

func formatYAMLError(data []byte, err error) error {
	line := extractYAMLErrorLine(err)
	if line == 0 {
		return err
	}

	entry := findYAMLPathAtLine(data, line)
	if entry.path == "" {
		return fmt.Errorf("yaml parse error at line %d: %w", line, err)
	}

	location := fmt.Sprintf("line %d", line)
	if entry.col > 0 {
		location = fmt.Sprintf("line %d, column %d", line, entry.col)
	}

	return fmt.Errorf("yaml parse error at %s (%s): %w", location, entry.path, err)
}

func extractYAMLErrorLine(err error) int {
	if err == nil {
		return 0
	}

	match := yamlLineRe.FindStringSubmatch(err.Error())
	if len(match) < 2 {
		return 0
	}

	line, parseErr := strconv.Atoi(match[1])
	if parseErr != nil {
		return 0
	}

	return line
}

func findYAMLPathAtLine(data []byte, line int) yamlPathEntry {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return yamlPathEntry{}
	}

	var entries []yamlPathEntry
	collectYAMLPaths(&root, nil, &entries)

	best := yamlPathEntry{}
	for i := range entries {
		entry := entries[i]
		if entry.line != line {
			continue
		}
		if best.path == "" || entry.depth > best.depth || (entry.depth == best.depth && entry.col < best.col) {
			best = entry
		}
	}

	return best
}

func collectYAMLPaths(node *yaml.Node, path []yamlPathPart, entries *[]yamlPathEntry) {
	if node == nil {
		return
	}

	recordYAMLPath(node, path, entries)

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			collectYAMLPaths(child, path, entries)
		}
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]
			childPath := appendYAMLKey(path, keyNode.Value)
			recordYAMLPath(keyNode, childPath, entries)
			collectYAMLPaths(valNode, childPath, entries)
		}
	case yaml.SequenceNode:
		for i, child := range node.Content {
			childPath := appendYAMLIndex(path, i)
			collectYAMLPaths(child, childPath, entries)
		}
	}
}

func recordYAMLPath(node *yaml.Node, path []yamlPathPart, entries *[]yamlPathEntry) {
	if node.Line == 0 {
		return
	}

	*entries = append(*entries, yamlPathEntry{
		line:  node.Line,
		col:   node.Column,
		depth: len(path),
		path:  formatYAMLPath(path),
	})
}

func appendYAMLKey(path []yamlPathPart, key string) []yamlPathPart {
	next := make([]yamlPathPart, len(path), len(path)+1)
	copy(next, path)
	next = append(next, yamlPathPart{kind: yamlPathKey, key: key})
	return next
}

func appendYAMLIndex(path []yamlPathPart, index int) []yamlPathPart {
	next := make([]yamlPathPart, len(path), len(path)+1)
	copy(next, path)
	next = append(next, yamlPathPart{kind: yamlPathIndex, index: index})
	return next
}

func formatYAMLPath(parts []yamlPathPart) string {
	if len(parts) == 0 {
		return ""
	}

	var b strings.Builder
	for _, part := range parts {
		switch part.kind {
		case yamlPathKey:
			if b.Len() == 0 {
				b.WriteString(part.key)
				continue
			}
			if isSimpleYAMLKey(part.key) {
				b.WriteString(".")
				b.WriteString(part.key)
				continue
			}
			b.WriteString("[\"")
			b.WriteString(escapeYAMLKey(part.key))
			b.WriteString("\"]")
		case yamlPathIndex:
			b.WriteString("[")
			b.WriteString(strconv.Itoa(part.index))
			b.WriteString("]")
		}
	}

	return b.String()
}

func isSimpleYAMLKey(key string) bool {
	if key == "" {
		return false
	}
	for i, r := range key {
		if r > 127 {
			return false
		}
		if i == 0 {
			if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && r != '_' {
				return false
			}
			continue
		}
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

func escapeYAMLKey(key string) string {
	key = strings.ReplaceAll(key, "\\", "\\\\")
	key = strings.ReplaceAll(key, "\"", "\\\"")
	return key
}
