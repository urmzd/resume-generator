package generators

import "github.com/urmzd/incipit/resume"

// TemplateData wraps a Resume with section ordering metadata for template rendering.
// Templates receive this as their dot context; the embedded *Resume makes all
// resume fields (Contact, Skills, etc.) accessible directly.
type TemplateData struct {
	*resume.Resume
	SectionOrder []string
}

// DefaultSectionOrder is the fallback when no order is provided.
var DefaultSectionOrder = []string{
	"summary", "certifications", "experience",
	"education", "skills", "projects", "languages",
}

// NewTemplateData creates a TemplateData from a resume and its section order.
// If sectionOrder is nil, DefaultSectionOrder is used.
func NewTemplateData(r *resume.Resume, sectionOrder []string) *TemplateData {
	if sectionOrder == nil {
		sectionOrder = DefaultSectionOrder
	}
	return &TemplateData{Resume: r, SectionOrder: sectionOrder}
}
