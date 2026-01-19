# Resume Schema Guide

## Overview

Resume Generator exposes a single JSON Schema for its unified resume format. Use it to:

- Validate resume data before generation
- Drive IDE autocompletion and inline errors
- Feed structured schema context into LLMs
- Integrate with external tools that consume JSON Schema

The schema covers the same fields used by the CLI for YAML, JSON, and TOML inputs.

## Resume Format (v2.0)

The schema reflects the `Resume` struct used at runtime. Key sections include:

```yaml
contact:
  name: string (required)
  email: string (required)
  phone: string
  location:
    city: string
    state: string
    country: string
  links:
    - uri: string

skills:
  title: string
  categories:
    - category: string
      items: [string]

experience:
  title: string
  positions:
    - company: string
      title: string
      highlights: [string]
      dates:
        start: datetime
        end: datetime
      location:
        city: string
        state: string
        country: string

projects:
  title: string
  projects:
    - name: string
      link:
        uri: string
      highlights: [string]

education:
  title: string
  institutions:
    - institution: string
      degree:
        name: string
        descriptions: [string]
      gpa:
        gpa: string
        max_gpa: string
      dates:
        start: datetime
        end: datetime
      location:
        city: string
        state: string
        country: string
      thesis:
        title: string
        highlights: [string]
        link:
          uri: string
```

## Generating the Schema

The CLI emits the schema to stdout by default:

```bash
resume-generator schema
```

To save it to a file:

```bash
resume-generator schema -o ./schemas/resume.schema.json
```

## Using the Schema with LLMs

1. Generate the schema:
   ```bash
   ./resume-generator schema -o ./schemas/resume.schema.json
   ```
2. Paste the schema into your LLM prompt alongside your request.

Example prompt:

```markdown
I'm working with a resume in the following format. Here's the JSON schema:

[Paste schema content]

Please help me:
1. Convert my plain-text resume into this format
2. Add metrics to my experience highlights
```

## IDE Integration (YAML Language Server)

Add a schema header to your resume file:

```yaml
# yaml-language-server: $schema=./schemas/resume.schema.json
```

Or configure VS Code settings:

```json
{
  "yaml.schemas": {
    "./schemas/resume.schema.json": ["resume*.yml", "resume*.yaml"]
  }
}
```

## Validation Examples

**Python (jsonschema)**:

```python
import json
from jsonschema import validate

with open('./schemas/resume.schema.json') as f:
    schema = json.load(f)

with open('resume.json') as f:
    resume = json.load(f)

validate(instance=resume, schema=schema)
```

**Node (ajv)**:

```js
const fs = require('fs');
const Ajv = require('ajv');

const schema = JSON.parse(fs.readFileSync('./schemas/resume.schema.json'));
const resume = JSON.parse(fs.readFileSync('resume.json'));

const ajv = new Ajv({ allErrors: true });
const validate = ajv.compile(schema);

if (!validate(resume)) {
  console.error(validate.errors);
}
```

## Troubleshooting

**Problem**: `resume-generator schema` prints an error

**Solution**: Make sure you're running the CLI from the repo root (or the binary is on your PATH) and that your Go build is up to date:

```bash
go build -o resume-generator .
./resume-generator schema
```
