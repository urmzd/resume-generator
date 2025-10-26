# Resume Schema Guide

## Overview

The resume-generator provides JSON Schema definitions for all supported resume formats. These schemas enable:

- 🤖 **LLM Integration** - Help AI models understand and generate resume content
- ✅ **Validation** - Verify resume data structure and types
- 📝 **IDE Support** - Enable autocompletion and inline validation
- 📚 **Documentation** - Auto-generate format documentation
- 🔧 **Tool Integration** - Easy integration with other tools

## Available Formats

### 1. Legacy Format (`resume-legacy.schema.json`)

**Best for**: Simple resumes, backward compatibility

**Features**:
- Simple, flat structure
- Easy to understand and edit
- Supports YAML, JSON, TOML
- Lightweight schema (~8KB)

**Key Fields**:
```yaml
contact:
  name: string (required)
  email: string (required)
  phone: string
  links: array of {link: string}

experience:
  - company: string
    title: string
    achievements: array of strings  # Job accomplishments
    dates:
      start: date
      end: date (optional, omit for current)
    location:
      city: string
      state: string

education:
  - school: string
    degree: string
    dates: {start, end}

skills:
  - category: string
    value: string (comma-separated)

projects:
  - name: string
    description: array of strings
```

### 2. Resume Format (`resume.schema.json`)

**Best for**: Detailed resumes, advanced features, customization

**Features**:
- Section ordering and visibility controls
- Quantifiable achievements with metrics
- Rich metadata and theming
- Multi-language support
- Comprehensive schema (~29KB)

**Key Fields**:
```yaml
meta:
  version: "2.0"
  output:
    formats: [pdf, html]
  theme: string

contact:
  order: int
  name: string (required)
  title: string
  email: string (required)
  phone: string
  links:
    - order: int
      text: string
      url: string
      type: string (github, linkedin, etc.)
  visibility:
    showPhone: bool
    showEmail: bool

experience:
  order: int
  title: string
  positions:
    - order: int
      company: string
      title: string
      highlights: array of strings
      achievements:
        - description: string
          metric:
            name: string
            before: number
            after: number
            unit: string
            improvement: number
      dates:
        start: datetime
        end: datetime (optional)
        current: bool

skills:
  categories:
    - name: string
      items:
        - name: string
          level: string
          years: int
```

### 3. JSON Resume Format (`resume-jsonresume.schema.json`)

**Best for**: Maximum compatibility, web hosting, standard format

**Features**:
- Community standard (jsonresume.org)
- Wide tool support
- Web-ready format
- Medium schema (~11KB)

**Key Fields**:
```json
{
  "basics": {
    "name": "string",
    "label": "string",
    "email": "string",
    "phone": "string",
    "url": "string",
    "summary": "string",
    "location": {
      "city": "string",
      "countryCode": "string",
      "region": "string"
    },
    "profiles": [
      {
        "network": "string",
        "username": "string",
        "url": "string"
      }
    ]
  },
  "work": [
    {
      "name": "string",
      "position": "string",
      "startDate": "date",
      "endDate": "date",
      "highlights": ["string"]
    }
  ],
  "education": [...],
  "skills": [...],
  "languages": [...],
  "interests": [...],
  "awards": [...],
  "publications": [...]
}
```

## Generating Schemas

### Command Line

```bash
# Generate all schemas
resume-generator schema generate -o ./schemas

# Generate specific format
resume-generator schema generate -f legacy -o ./schemas
resume-generator schema generate -f resume -o ./schemas

# List available formats
resume-generator schema list

# Print schema to stdout (for piping)
resume-generator schema generate -f legacy
```

## Using Schemas with LLMs

### 1. Claude / ChatGPT / GPT-4

**Prompt Template**:

```markdown
I'm working with a resume in the following format. Here's the JSON schema:

[Paste schema content]

Based on this schema, please help me:
1. [Your specific task]

Example tasks:
- "Convert my plain text resume to this format"
- "Enhance my existing resume data"
- "Generate a sample resume for a software engineer"
- "Add quantifiable metrics to my achievements"
- "Restructure my experience section"
```

**Example Usage**:

```bash
# Get schema
./resume-generator schema generate -f resume > schema.json

# Use in LLM prompt
cat schema.json | pbcopy  # Copy to clipboard (macOS)
```

Then paste into your LLM chat with your request.

### 2. LangChain / LlamaIndex

```python
import json
from langchain.output_parsers import StructuredOutputParser
from langchain.prompts import PromptTemplate

# Load schema
with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)

# Create structured parser
parser = StructuredOutputParser.from_response_schemas([
    # Convert JSON schema to LangChain schema
])

# Use in chain
prompt = PromptTemplate(
    template="Convert this text to a resume:\n{input}\n\n{format_instructions}",
    input_variables=["input"],
    partial_variables={"format_instructions": parser.get_format_instructions()}
)
```

### 3. Structured Output APIs (OpenAI, Anthropic)

**OpenAI Function Calling**:

```python
import openai
import json

# Load schema
with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)

response = openai.ChatCompletion.create(
    model="gpt-4",
    messages=[
        {"role": "system", "content": "You are a resume formatting assistant."},
        {"role": "user", "content": "Convert this text to a structured resume: ..."}
    ],
    functions=[{
        "name": "create_resume",
        "description": "Create a structured resume",
        "parameters": schema
    }],
    function_call={"name": "create_resume"}
)

resume_data = json.loads(response.choices[0].message.function_call.arguments)
```

**Anthropic Claude with Tool Use**:

```python
import anthropic
import json

client = anthropic.Anthropic()

with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)

response = client.messages.create(
    model="claude-3-sonnet-20240229",
    max_tokens=4096,
    tools=[{
        "name": "create_resume",
        "description": "Create a structured resume",
        "input_schema": schema
    }],
    messages=[{
        "role": "user",
        "content": "Convert this text to a structured resume: ..."
    }]
)
```

## Validation

### 1. Using JSON Schema Validators

**Python (jsonschema)**:

```python
import json
import yaml
from jsonschema import validate, ValidationError

# Load schema
with open('assets/schema/resume-legacy.schema.json') as f:
    schema = json.load(f)

# Load resume
with open('resume.yml') as f:
    resume = yaml.safe_load(f)

# Validate
try:
    validate(instance=resume, schema=schema)
    print("✓ Resume is valid!")
except ValidationError as e:
    print(f"✗ Validation error: {e.message}")
```

**Node.js (ajv)**:

```javascript
const Ajv = require('ajv');
const fs = require('fs');
const yaml = require('js-yaml');

const ajv = new Ajv();

// Load schema
const schema = JSON.parse(fs.readFileSync('assets/schema/resume-legacy.schema.json'));

// Load resume
const resume = yaml.load(fs.readFileSync('resume.yml'));

// Validate
const valid = ajv.validate(schema, resume);
if (valid) {
    console.log('✓ Resume is valid!');
} else {
    console.log('✗ Validation errors:', ajv.errors);
}
```

### 2. Using the CLI

```bash
# The tool automatically validates on generation
./resume-generator run -i resume.yml -o output-directory

# Explicit validation
./resume-generator validate resume.yml
```

## IDE Integration

### VSCode

1. Install the [YAML extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)

2. Add to your `resume.yml`:

```yaml
# yaml-language-server: $schema=./assets/schema/resume-legacy.schema.json

contact:
  name: John Doe  # Now with autocomplete!
```

3. Or configure globally in `.vscode/settings.json`:

```json
{
  "yaml.schemas": {
    "./assets/schema/resume-legacy.schema.json": ["resume*.yml", "resume*.yaml"]
  }
}
```

### JetBrains IDEs (IntelliJ, PyCharm, etc.)

1. Open Settings → Languages & Frameworks → Schemas and DTDs → JSON Schema Mappings

2. Add mapping:
   - Schema file: `assets/schema/resume-legacy.schema.json`
   - File pattern: `resume*.yml`

## Advanced LLM Prompting Techniques

### 1. Resume Enhancement

```markdown
Given this resume schema: [schema]

And this existing resume data: [current resume]

Please enhance the resume by:
1. Adding quantifiable metrics where possible (e.g., "increased by X%")
2. Using action verbs for achievements
3. Ensuring consistency in date formats
4. Adding missing but inferable information
5. Improving clarity and impact of descriptions

Return the resume resume in the same format.
```

### 2. Format Conversion

```markdown
Source Format: [paste legacy schema]
Target Format: [paste resume schema]

Current Resume: [paste current data]

Please convert this resume from the legacy format to the resume format,
preserving all information and adding appropriate default values for new
fields (order: sequential, visibility: all true).
```

### 3. Resume Generation from Text

```markdown
Schema: [paste schema]

Please convert this unstructured text into a properly formatted resume
following the schema above:

[Paste plain text resume or LinkedIn profile]

Ensure:
- Proper date parsing (YYYY-MM-DD format)
- Consistent formatting
- Logical section ordering
- No information loss
```

### 4. Multi-Language Resume

```markdown
Schema: [paste resume schema]
Current Resume: [English resume]

Please create a Spanish version of this resume, translating all fields
while maintaining the same structure. Update the meta.language field to "es".
```

## Schema Customization

You can extend the schemas for custom use cases:

```bash
# Generate base schema
./resume-generator schema generate -f resume -o ./schemas

# Edit and add custom fields
# Add to your resume with custom fields
# The tool will pass through unknown fields in some formats
```

## Best Practices

### For LLM Integration

1. **Always include schema** - Don't rely on LLMs "knowing" the format
2. **Provide examples** - Show example resumes in your prompts
3. **Validate output** - Always validate LLM-generated resumes
4. **Iterate** - Use multi-step prompts for better results:
   - Step 1: Extract raw data
   - Step 2: Structure data
   - Step 3: Enhance and optimize

### For Validation

1. **Validate early** - Check structure before attempting generation
2. **Use strict mode** - Enable all schema validations
3. **Provide context** - Include field descriptions in error messages

### For IDE Integration

1. **Use relative paths** - Makes schemas portable across machines
2. **Commit schemas** - Include in version control
3. **Update regularly** - Regenerate when tool is updated

## Troubleshooting

### Schema Generation Issues

**Problem**: "Command not found: schema"
**Solution**: Rebuild the tool: `go build -o resume-generator`

**Problem**: Schema missing fields
**Solution**: Regenerate schemas after updating to latest version

### Validation Issues

**Problem**: Valid resume fails validation
**Solution**: Check that you're using the matching schema version

**Problem**: Date format errors
**Solution**: Use ISO 8601 format: `YYYY-MM-DD` or `YYYY-MM-DDTHH:MM:SSZ`

### LLM Integration Issues

**Problem**: LLM returns invalid JSON
**Solution**: Use structured output APIs or add validation step

**Problem**: Schema too large for context window
**Solution**: Use the simpler "legacy" format or extract relevant portions

## Examples

### Complete LLM Workflow

```bash
# 1. Generate schema
./resume-generator schema generate -f legacy -o ./schemas

# 2. Use with LLM to create resume (copy schema + paste in chat)
# LLM outputs: resume-data.yml

# 3. Validate
./resume-generator validate resume-data.yml

# 4. Generate PDF
./resume-generator run -i resume-data.yml -o resume-exports -t modern-latex

# 5. If validation fails, fix and repeat
```

### Python Automation Script

```python
#!/usr/bin/env python3
import json
import yaml
import subprocess
from anthropic import Anthropic

# Generate schema
subprocess.run([
    './resume-generator', 'schema', 'generate',
    '-f', 'legacy', '-o', './schemas'
])

# Load schema
with open('./schemas/resume-legacy.schema.json') as f:
    schema = json.load(f)

# Use LLM to generate resume
client = Anthropic()
response = client.messages.create(
    model="claude-3-sonnet-20240229",
    max_tokens=4096,
    tools=[{
        "name": "create_resume",
        "description": "Create a structured resume",
        "input_schema": schema
    }],
    messages=[{
        "role": "user",
        "content": "Create a resume for a senior software engineer with 10 years of experience..."
    }]
)

# Extract and save resume
resume_data = response.content[0].input
with open('generated-resume.yml', 'w') as f:
    yaml.dump(resume_data, f)

# Generate PDF
subprocess.run([
    './resume-generator', 'run',
    '-i', 'generated-resume.yml',
    '-o', 'generated-resume.pdf',
    '-t', 'modern-latex'
])

print("✓ Resume generated: generated-resume.pdf")
```

## Resources

- [JSON Schema Documentation](https://json-schema.org/)
- [JSON Resume Standard](https://jsonresume.org/)
- [Resume Generator Examples](./assets/example_inputs/)
- [Template Documentation](./README.md#templates)

## Support

For issues or questions:
- GitHub Issues: https://github.com/urmzd/resume-generator/issues
- Schema Validation: Use `resume-generator validate`
- Format Questions: Use `resume-generator schema list`
