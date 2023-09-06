# Generate Resumes

Generate beautiful $\LaTeX$ resumes from a single source using a single command.

## Usage

```bash
  ./generate-resumes --help
  # we need to ensure the output folder exists
  mkdir output && ./generate-resumes examples/config.toml -o "$(pwd)/output"
```

## Requirements

- `xelatex` (requires both the `enumitem` and `titlesec` packages)
