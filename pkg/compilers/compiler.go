package compilers

// Compiler interface for LaTeX compilation.
// Used by compilers package to compile LaTeX documents to PDF.
type Compiler interface {
	LoadClasses(string)
	AddOutputFolder(string)
	Compile(string, string) string
}
