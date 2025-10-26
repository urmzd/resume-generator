package definition

// Compiler interface for LaTeX compilation.
// Used by compilers package to compile LaTeX documents to PDF.
type Compiler interface {
	LoadClasses(string)
	AddOutputFolder(string)
	Compile(string, string) string
}

// Location represents a physical location with various levels of detail
type Location struct {
	Address     string `json:"address,omitempty" yaml:"address,omitempty" toml:"address,omitempty"`
	PostalCode  string `json:"postalCode,omitempty" yaml:"postalCode,omitempty" toml:"postalCode,omitempty"`
	City        string `json:"city" yaml:"city" toml:"city"`
	State       string `json:"state,omitempty" yaml:"state,omitempty" toml:"state,omitempty"`
	Country     string `json:"country,omitempty" yaml:"country,omitempty" toml:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty" yaml:"countryCode,omitempty" toml:"countryCode,omitempty"`
	Region      string `json:"region,omitempty" yaml:"region,omitempty" toml:"region,omitempty"`
}

// CategoryValuePair represents a key-value pair for structured data
type CategoryValuePair struct {
	Category string `json:"category" yaml:"category" toml:"category"`
	Value    string `json:"value" yaml:"value" toml:"value"`
}
