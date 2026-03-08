package resume

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDegree_UnmarshalYAML_Note(t *testing.T) {
	tests := []struct {
		name             string
		yaml             string
		wantName         string
		wantDescriptions []string
	}{
		{
			name:             "standard degree with descriptions",
			yaml:             "name: BSc\ndescriptions:\n  - Dean's List\n  - Honors",
			wantName:         "BSc",
			wantDescriptions: []string{"Dean's List", "Honors"},
		},
		{
			name:             "note as alias for descriptions",
			yaml:             "name: MSW\nnote: Clinical focus",
			wantName:         "MSW",
			wantDescriptions: []string{"Clinical focus"},
		},
		{
			name:             "descriptions takes precedence over note",
			yaml:             "name: BSc\ndescriptions:\n  - From descriptions\nnote: From note",
			wantName:         "BSc",
			wantDescriptions: []string{"From descriptions"},
		},
		{
			name:             "no descriptions or note",
			yaml:             "name: PhD",
			wantName:         "PhD",
			wantDescriptions: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Degree
			if err := yaml.Unmarshal([]byte(tt.yaml), &d); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if d.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", d.Name, tt.wantName)
			}
			if len(d.Descriptions) != len(tt.wantDescriptions) {
				t.Fatalf("Descriptions len = %d, want %d", len(d.Descriptions), len(tt.wantDescriptions))
			}
			for i, desc := range d.Descriptions {
				if desc != tt.wantDescriptions[i] {
					t.Errorf("Descriptions[%d] = %q, want %q", i, desc, tt.wantDescriptions[i])
				}
			}
		})
	}
}

func TestEducation_UnmarshalYAML_Credential(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		wantDegree string
		wantDescs  []string
	}{
		{
			name:       "standard degree field",
			yaml:       "institution: MIT\ndegree:\n  name: BSc\n  descriptions:\n    - Honors\ndates:\n  start: 2020-01-01",
			wantDegree: "BSc",
			wantDescs:  []string{"Honors"},
		},
		{
			name:       "credential as alias for degree",
			yaml:       "institution: UofT\ncredential:\n  name: MSW\n  note: Clinical focus\ndates:\n  start: 2020-01-01",
			wantDegree: "MSW",
			wantDescs:  []string{"Clinical focus"},
		},
		{
			name:       "degree takes precedence over credential",
			yaml:       "institution: MIT\ndegree:\n  name: BSc\ncredential:\n  name: MSW\ndates:\n  start: 2020-01-01",
			wantDegree: "BSc",
			wantDescs:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Education
			if err := yaml.Unmarshal([]byte(tt.yaml), &e); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if e.Degree.Name != tt.wantDegree {
				t.Errorf("Degree.Name = %q, want %q", e.Degree.Name, tt.wantDegree)
			}
			if len(e.Degree.Descriptions) != len(tt.wantDescs) {
				t.Fatalf("Descriptions len = %d, want %d", len(e.Degree.Descriptions), len(tt.wantDescs))
			}
			for i, desc := range e.Degree.Descriptions {
				if desc != tt.wantDescs[i] {
					t.Errorf("Descriptions[%d] = %q, want %q", i, desc, tt.wantDescs[i])
				}
			}
		})
	}
}
