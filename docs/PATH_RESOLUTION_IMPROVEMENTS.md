# Path Resolution Improvements

## Overview

This document describes the comprehensive path resolution improvements made to the resume-generator CLI tool. These changes allow the tool to support flexible input/output paths, work from any directory, and handle various path formats.

## Problems Solved

### Before

1. **Hardcoded output filename**: Always output to `resume.pdf`
2. **Hardcoded asset paths**: Templates and LaTeX classes assumed specific working directory
3. **Working directory dependency**: Failed when run from different directories
4. **Unused CLI variables**: Variables defined but never utilized
5. **No path validation**: No early validation of input/output paths
6. **No home directory support**: Couldn't use `~` for home directory paths

### After

1. **Flexible output workspaces**: Custom output paths or auto-generated timestamped directories per run
2. **Dynamic asset resolution**: Assets found relative to executable or working directory
3. **Directory-independent operation**: Works from any directory
4. **Proper CLI flag usage**: All defined flags now functional
5. **Early path validation**: Input files validated before processing
6. **Full path support**: Relative, absolute, and home directory paths

## New Features

### 1. Path Utility Package (`pkg/utils/paths.go`)

Created a centralized path resolution utility with the following functions:

- **`ResolvePath(path string)`**: Resolves any path to absolute form
  - Supports relative paths (resolved from cwd)
  - Supports absolute paths (cleaned and normalized)
  - Supports home directory expansion (`~/path`)

- **`ResolveOutputPath(path string, createDir bool)`**: Resolves output paths
  - Creates parent directories if needed
  - Handles both file and directory paths

- **`ResolveAssetPath(relativePath string)`**: Multi-location template/resource resolution
  - Tries path relative to current working directory
  - Tries path relative to executable directory
  - Ensures templates work in development and deployed scenarios

- **`GetExecutableDir()`**: Gets executable directory
  - Resolves symlinks
  - Useful for finding bundled templates

- **Helper functions**:
  - `EnsureDir()`: Creates directory if missing
  - `FileExists()`: Check if file exists
  - `DirExists()`: Check if directory exists

### 2. Per-Run Output Workspace

New directory format: `{first}_{optional_middle}_{last}_{iso8601}/`

Structure:
- `resume.pdf` (or custom name when `--output` specifies a file) saved at the root of the run directory
- `debug/` subfolder containing generated `.tex`, `.html`, `.log`, `.aux`, and any LaTeX support files copied from the template

Examples:
- `john_doe_2025-10-25T13-45-30/resume.pdf`
- `jane_marie_smith_2025-10-25T13-45-30/debug/resume.tex`

Benefits:
- Sanitizes names (removes special characters, spaces)
- Includes timestamp for versioning and collision avoidance
- Keeps all intermediate artifacts together for easier debugging

### 3. Enhanced CLI Flags

#### Run Command (`resume-generator run`)

**New/Updated Flags:**
- `-o, --output <path>`: Specify exact output file path
  - Supports any path format (relative, absolute, home)
  - If path is directory, creates a per-run subfolder with sanitized name + timestamp
  - Creates parent directories automatically

- `--output-dir <path>`: Deprecated but still supported for backwards compatibility

**Examples:**
```bash
# Custom output file
./resume-generator run -i resume.yml -o ~/Documents/my_resume.pdf

# Output to directory (creates timestamped folder)
./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o output/
# Creates: output/john_doe_2025-10-25T13-45-30/resume.pdf

# Relative paths
./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o output/resume.pdf
```

### 4. Path Resolution in All Commands

All commands now properly resolve paths:

- **`run`**: Input file and output path handling
- **`preview`**: Input file with validation
- **`validate`**: Input file with early existence check
- **`templates validate`**: Template file resolution

### 5. Template Path Resolution

Templates are now found using intelligent multi-location search:

1. **Working Directory**: `./templates/`
2. **Executable Directory**: `{exe_dir}/templates/`
3. **Fallback**: Returns best-effort path for error reporting

This allows the tool to:
- Work in development (from repo root)
- Work when installed (from any directory)
- Work in containers (with template directories mounted)

## Changes by File

### New Files

- `pkg/utils/paths.go`: Centralized path resolution utilities

### Modified Files

#### `cmd/root.go`
- No changes (already had variable definitions)

#### `cmd/run.go`
- Added path utilities import
- Added `time` import for timestamps
- Updated CLI flags to use `OutputFile`
- Added input path resolution and validation
- Added smart output path determination
- Added per-run output directory creation with timestamped naming
- Updated LaTeX compilation to copy template-local support files and preserve debug artifacts

#### `cmd/preview.go`
- Added path utilities import
- Added input path resolution and validation

#### `cmd/validate.go`
- Added path utilities import
- Added input path resolution and validation

#### `cmd/templates.go`
- Added path utilities import
- Added template path resolution in validate subcommand

#### `pkg/generators/generator.go`
- Added path utilities import
- Updated `LoadTemplate()` to use `ResolveAssetPath()`
- Updated `ListTemplates()` to use `ResolveAssetPath()`
- Added better error messages with actual search paths

## Testing

### Test Scenarios

1. **Relative Paths**
   ```bash
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o output.pdf
   ```

2. **Absolute Paths**
   ```bash
   ./resume-generator run -i /full/path/to/resume.yml -o /full/path/output.pdf
   ```

3. **Home Directory**
   ```bash
   ./resume-generator run -i ~/resumes/resume.yml -o ~/Documents/resume.pdf
   ```

4. **Directory as Output**
   ```bash
   ./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o output/
   # Creates: output/john_doe_2025-10-25T13-45-30/resume.pdf
   ```

5. **Different Working Directories**
   ```bash
   cd /tmp
   /path/to/resume-generator templates list
   # Still finds templates correctly
   ```

### Verified Functionality

✅ Templates list works from any directory
✅ Preview resolves relative input paths
✅ Validate resolves relative input paths
✅ Run command resolves input/output paths
✅ Assets found relative to executable
✅ Home directory expansion works
✅ Per-run workspace creation works
✅ Output directory creation works
✅ Template-local LaTeX support files copied automatically

## Backwards Compatibility

All changes are backwards compatible:

- **`--output-dir` flag**: Still works (deprecated but functional)
- **Default behavior**: If no output specified, a timestamped workspace directory is created automatically
- **Existing scripts**: Will continue to work with new functionality

## Migration Guide

### For Users

**Old way:**
```bash
./resume-generator run -i resume.yml
# Output: ./john_doe_2025-10-25T13-45-30/resume.pdf
```

**New way (same result):**
```bash
./resume-generator run -i resume.yml
# Output: ./john_doe_2025-10-25T13-45-30/resume.pdf
```

**New capabilities:**
```bash
# Custom output location
./resume-generator run -i resume.yml -o ~/Documents/my_resume.pdf

# Custom output directory with timestamped workspace
./resume-generator run -i assets/example_inputs/sample-enhanced.yml -o output/
# Creates: output/john_doe_2025-10-25T13-45-30/resume.pdf
```

### For Developers

**Resolving paths in code:**
```go
// Old way (don't do this)
filepath.Join("templates", name)

// New way
utils.ResolveAssetPath(filepath.Join("templates", name))
```

**Validating input files:**
```go
// Old way
data, err := os.ReadFile(inputPath)

// New way
resolved, err := utils.ResolvePath(inputPath)
if err != nil { /* handle */ }
if !utils.FileExists(resolved) { /* handle */ }
data, err := os.ReadFile(resolved)
```

## Future Enhancements

Potential improvements for future versions:

1. **Environment Variable Support**: `RESUME_TEMPLATES_DIR`, `RESUME_OUTPUT_DIR`
2. **Config File**: `.resume-generator.yaml` for default paths
3. **Template Discovery**: Search multiple template directories
4. **Symlink Resolution**: Better handling of symlinked template directories
5. **URL Support**: Download templates from URLs
6. **Relative Output**: Output relative to input file location

## Summary

These improvements make the resume-generator CLI tool significantly more flexible and user-friendly. Users can now:

- Use any path format they prefer
- Run the tool from any directory
- Get descriptive output filenames automatically
- Customize template locations
- Trust that paths will be validated early

The changes maintain full backwards compatibility while adding powerful new capabilities for path handling.
