# LaTeX Engine Support

## Overview

The resume-generator now supports multiple LaTeX compilation engines, providing flexibility and compatibility across different LaTeX installations.

## Supported Engines

The tool supports the following LaTeX engines (in order of preference):

1. **xelatex** - Modern engine with excellent Unicode and font support
2. **pdflatex** - Traditional and widely available engine
3. **lualatex** - Modern engine with Lua scripting capabilities
4. **latex** - Classic LaTeX engine (requires DVI to PDF conversion)

## Auto-Detection

By default, the tool automatically detects the first available LaTeX engine on your system and uses it for compilation. The detection order prioritizes modern engines with better Unicode support.

### How it Works

When you run a LaTeX compilation without specifying an engine:

```bash
./resume-generator run -i resume.yml -t base-latex
```

The tool will:
1. Check for `xelatex` in your system PATH
2. If not found, check for `pdflatex`
3. If not found, check for `lualatex`
4. If not found, check for `latex`
5. Use the first one found, or show a helpful error message if none are available

## Manual Engine Selection

You can explicitly specify which LaTeX engine to use with the `--latex-engine` or `-e` flag:

```bash
# Use pdflatex specifically
./resume-generator run -i resume.yml -t base-latex --latex-engine pdflatex

# Use xelatex specifically
./resume-generator run -i resume.yml -t base-latex -e xelatex

# Use lualatex specifically
./resume-generator run -i resume.yml -t base-latex -e lualatex
```

## Checking Available Engines

To see which LaTeX engines are installed on your system:

```bash
./resume-generator templates engines
```

Example output:

```
Checking for LaTeX engines...

✓ Found 2 LaTeX engine(s):

✓ xelatex (default - will be used if no engine is specified)
  pdflatex

Usage:
  # Use default engine (auto-detected)
  resume-generator run -i resume.yml -t base-latex

  # Specify a particular engine
  resume-generator run -i resume.yml -t base-latex --latex-engine xelatex
```

## Installing LaTeX

If no LaTeX engines are found, you'll receive helpful installation instructions.

### macOS

```bash
# Install MacTeX (includes all engines)
brew install --cask mactex

# Or install BasicTeX (smaller, core engines only)
brew install --cask basictex
```

### Linux

```bash
# Debian/Ubuntu
sudo apt-get install texlive-xetex texlive-latex-base

# Fedora/RHEL
sudo dnf install texlive-xetex texlive-latex

# Arch Linux
sudo pacman -S texlive-core texlive-latexextra
```

### Windows

Download and install one of:
- [TeX Live](https://www.tug.org/texlive/) - Full-featured distribution
- [MiKTeX](https://miktex.org/) - User-friendly distribution with package manager

### Docker

Use a Docker image with LaTeX pre-installed:

```bash
docker run --rm -v $(pwd):/work texlive/texlive \
  resume-generator run -i /work/resume.yml -o /work/output.pdf -t base-latex
```

## Engine-Specific Considerations

### xelatex
- **Best for**: Modern documents, Unicode text, custom fonts
- **Pros**: Excellent font support, handles UTF-8 natively
- **Cons**: Slightly slower than pdflatex

### pdflatex
- **Best for**: Traditional LaTeX documents, compatibility
- **Pros**: Fast compilation, widely available
- **Cons**: Limited Unicode support, requires font packages

### lualatex
- **Best for**: Complex documents, programmable features
- **Pros**: Lua scripting, good font support
- **Cons**: Less common than xelatex/pdflatex

### latex
- **Best for**: Legacy documents
- **Pros**: Extremely compatible
- **Cons**: Produces DVI (needs conversion to PDF)

## Template Compatibility

All included templates are designed to work with any of the supported engines. However:

- Templates using custom fonts work best with **xelatex** or **lualatex**
- Templates with special characters require **xelatex**, **lualatex**, or proper encoding setup in **pdflatex**
- Simple templates work well with any engine

## Troubleshooting

### Engine Not Found

**Error:**
```
no LaTeX engine found

Please install one of the following:
  - TeX Live:   https://www.tug.org/texlive/
  - MiKTeX:     https://miktex.org/
  - MacTeX:     https://www.tug.org/mactex/ (macOS)
```

**Solution:** Install a LaTeX distribution (see "Installing LaTeX" above)

### Specific Engine Not Working

**Error:**
```
LaTeX compilation error with pdflatex: exit status 1
```

**Solution:**
1. Check if the engine is properly installed: `which pdflatex`
2. Try a different engine: `--latex-engine xelatex`
3. Check the saved `.tex` file for syntax errors

### Font Errors with pdflatex

**Error:**
```
! Font \U/fontname/m/n/10=fontfile at 10.0pt not loadable
```

**Solution:** Use xelatex or lualatex for better font support:
```bash
./resume-generator run -i resume.yml -t base-latex -e xelatex
```

## Advanced Usage

### Using with Different Templates

```bash
# Auto-detect engine for any template
./resume-generator run -i resume.yml -t base-latex
./resume-generator run -i resume.yml -t json-resume-latex

# Force specific engine for all templates
./resume-generator run -i resume.yml -t base-latex -e pdflatex
```

### Debugging LaTeX Output

The tool saves the generated LaTeX source alongside the PDF:

```bash
./resume-generator run -i resume.yml -o output.pdf -t base-latex

# This creates:
# - output.pdf (compiled PDF)
# - output.tex (LaTeX source for debugging)
```

You can compile the `.tex` file manually to see detailed errors:

```bash
xelatex output.tex
# or
pdflatex output.tex
```

## API Changes

### For Developers

If you're integrating the resume-generator as a library:

```go
import "github.com/urmzd/resume-generator/pkg/compilers"

// Auto-detect engine
compiler, err := compilers.NewAutoLaTeXCompiler(logger)
if err != nil {
    // No LaTeX engine found
}

// Or specify engine
compiler := compilers.NewLaTeXCompiler("xelatex", logger)

// Check available engines
engines := compilers.GetAvailableLaTeXEngines()
// Returns: []string{"xelatex", "pdflatex", ...}

// Detect default engine
engine := compilers.DetectLaTeXEngine()
// Returns: "xelatex" or "" if none found
```

## Summary

The multi-engine LaTeX support provides:

✅ **Automatic detection** - Works out of the box with any installed engine
✅ **Manual override** - Specify engine when needed
✅ **Better error messages** - Clear guidance when LaTeX is not installed
✅ **Flexibility** - Use the best engine for your needs
✅ **Compatibility** - Works with all major LaTeX distributions

For most users, the auto-detection feature means you can simply run:

```bash
./resume-generator run -i resume.yml -t base-latex
```

And the tool will find and use the best available LaTeX engine on your system!
