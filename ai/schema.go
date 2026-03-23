package ai

import (
	"github.com/invopop/jsonschema"
	"github.com/urmzd/incipit/resume"
	"github.com/urmzd/saige/agent/types"
)

// ResumeSchema builds a saige ParameterSchema from the Resume struct,
// suitable for use as ResponseSchema in structured output.
func ResumeSchema() *types.ParameterSchema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	schema := reflector.Reflect(&resume.Resume{})
	ps := convertSchema(schema)
	return &ps
}

// convertSchema converts a jsonschema.Schema to a saige ParameterSchema.
func convertSchema(s *jsonschema.Schema) types.ParameterSchema {
	ps := types.ParameterSchema{
		Type:       s.Type,
		Properties: make(map[string]types.PropertyDef),
	}

	if s.Required != nil {
		ps.Required = s.Required
	}

	if s.Properties != nil {
		for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
			ps.Properties[pair.Key] = convertPropertyDef(pair.Value)
		}
	}

	return ps
}

// convertPropertyDef converts a jsonschema.Schema (used as a property) to a PropertyDef.
func convertPropertyDef(s *jsonschema.Schema) types.PropertyDef {
	pd := types.PropertyDef{
		Type:        s.Type,
		Description: s.Description,
	}

	if len(s.Enum) > 0 {
		for _, e := range s.Enum {
			if str, ok := e.(string); ok {
				pd.Enum = append(pd.Enum, str)
			}
		}
	}

	if s.Items != nil {
		item := convertPropertyDef(s.Items)
		pd.Items = &item
	}

	if s.Properties != nil {
		pd.Properties = make(map[string]types.PropertyDef)
		for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
			pd.Properties[pair.Key] = convertPropertyDef(pair.Value)
		}
		pd.Required = s.Required
	}

	return pd
}
