# Resume JSON Schemas

This directory contains JSON Schema definitions for all supported resume formats.

## Files

- **`resume-legacy.schema.json`** (7.9KB) - Simple, traditional resume format
- **`resume-enhanced.schema.json`** (29KB) - Advanced format with ordering, metrics, visibility
- **`resume-jsonresume.schema.json`** (11KB) - Standard JSON Resume format (jsonresume.org)

## Quick Start

### Generate Schemas

```bash
# Regenerate all schemas
resume-generator schema generate -o assets/schema

# Generate specific format
resume-generator schema generate -f legacy -o assets/schema
```

### Use in IDE

**VSCode** - Add to your resume YAML:

```yaml
# yaml-language-server: $schema=assets/schema/resume-legacy.schema.json

contact:
  name: Your Name
  # IDE now provides autocomplete!
```

### Use with LLMs

```bash
# Get schema for Claude/GPT
cat assets/schema/resume-legacy.schema.json | pbcopy

# Or in your code
import json
with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)
    # Use with LLM API
```

### Validate Resume

```python
import json
import yaml
from jsonschema import validate

# Load schema
with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)

# Load resume
with open('resume.yml') as f:
    resume = yaml.safe_load(f)

# Validate
validate(instance=resume, schema=schema)
print("✓ Valid!")
```

## Schema Features

Each schema includes:

- ✅ **Complete type definitions** - All fields with proper types
- 📝 **Field descriptions** - Explain what each field means
- 🔍 **Validation rules** - Required fields, formats, constraints
- 📋 **Examples** - Sample resumes in each format
- 🎯 **Metadata** - Schema version, title, description

## Format Comparison

| Feature | Legacy | Enhanced | JSON Resume |
|---------|--------|----------|-------------|
| Simplicity | ★★★★★ | ★★☆☆☆ | ★★★★☆ |
| Features | ★★☆☆☆ | ★★★★★ | ★★★★☆ |
| Compatibility | ★★★★★ | ★★☆☆☆ | ★★★★★ |
| Schema Size | 7.9KB | 29KB | 11KB |
| Section Ordering | ❌ | ✅ | ❌ |
| Visibility Controls | ❌ | ✅ | ❌ |
| Metrics | ❌ | ✅ | ❌ |
| LLM Friendly | ✅ | ⚠️ Complex | ✅ |

## Example: Using with Claude

```python
import anthropic
import json

client = anthropic.Anthropic()

# Load schema
with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)

# Generate resume with structured output
response = client.messages.create(
    model="claude-3-sonnet-20240229",
    max_tokens=4096,
    tools=[{
        "name": "create_resume",
        "description": "Create a structured resume from text",
        "input_schema": schema
    }],
    messages=[{
        "role": "user",
        "content": """
        Convert this to a structured resume:

        John Doe
        Software Engineer
        john@example.com

        Experience:
        - Tech Corp (2020-2023): Built microservices, reduced latency by 40%
        - Startup Inc (2018-2020): Led team of 5, launched 3 products

        Skills: Python, Go, AWS, Kubernetes
        """
    }]
)

# Extract resume data
resume = response.content[0].input
print(resume)
```

## More Information

See [SCHEMA_GUIDE.md](../../docs/SCHEMA_GUIDE.md) for complete documentation on:

- Using schemas with different LLMs
- Validation techniques
- IDE integration
- Best practices
- Troubleshooting
- Advanced examples

## Regenerating Schemas

Schemas should be regenerated when:

- Tool is updated to new version
- Data structures change
- New fields are added

```bash
# Always generate from project root
cd /path/to/resume-generator
./resume-generator schema generate -o assets/schema
```

## Schema Standards

These schemas follow:

- **JSON Schema Draft 2020-12** - Latest JSON Schema standard
- **OpenAPI 3.0 compatible** - Can be used in OpenAPI specs
- **Strict validation** - `additionalProperties: false` for type safety
