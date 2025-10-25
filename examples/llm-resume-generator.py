#!/usr/bin/env python3
"""
LLM Resume Generator Example using Ollama

This script demonstrates how to use the resume schemas with Ollama
to generate structured resume data from unstructured text.

Requirements:
    pip install pyyaml jsonschema requests
    # Install Ollama: https://ollama.ai

Usage:
    # Start Ollama and pull a model
    ollama pull llama3.2

    # Run the script
    python examples/llm-resume-generator.py
"""

import json
import yaml
import requests
from pathlib import Path


def load_schema(format_name='legacy'):
    """Load a resume schema file."""
    schema_path = Path(__file__).parent / 'schemas' / f'resume-{format_name}.schema.json'
    with open(schema_path) as f:
        return json.load(f)


def generate_with_ollama(text_resume: str, format_name='legacy', model='llama3.2'):
    """
    Generate structured resume using Ollama (local LLM).

    Args:
        text_resume: Unstructured resume text
        format_name: 'legacy', 'enhanced', or 'json-resume'
        model: Ollama model name (llama3.2, mistral, etc.)

    Returns:
        dict: Structured resume data
    """
    schema = load_schema(format_name)

    # Get example from schema for reference
    example = None
    if 'examples' in schema and schema['examples']:
        example = schema['examples'][0]

    print(f"🤖 Generating resume with Ollama {model} (format: {format_name})...")

    # Build prompt with schema information
    prompt = f"""You are a resume formatting expert. Convert the following unstructured resume text into a structured JSON format.

OUTPUT FORMAT:
{json.dumps(example, indent=2) if example else "See schema below"}

SCHEMA REQUIREMENTS:
- contact.name: required string
- contact.email: required string
- contact.phone: string
- experience: array of objects with company, title, achievements (array of strings), dates
- education: array with school, degree, dates
- skills: array with category and value (comma-separated string)
- Dates must be in YYYY-MM-DD format
- Use "achievements" field for bullet points

RESUME TEXT TO CONVERT:
{text_resume}

Return ONLY valid JSON matching the schema. No explanations, just the JSON object.
"""

    try:
        # Call Ollama API
        response = requests.post(
            'http://localhost:11434/api/generate',
            json={
                'model': model,
                'prompt': prompt,
                'stream': False,
                'format': 'json',  # Request JSON output
                'options': {
                    'temperature': 0.1,  # Lower temperature for more consistent output
                }
            },
            timeout=120
        )

        response.raise_for_status()
        result = response.json()

        # Extract the generated JSON
        generated_text = result['response']

        # Parse JSON from response
        try:
            resume_data = json.loads(generated_text)
            return resume_data
        except json.JSONDecodeError as e:
            print(f"❌ Failed to parse JSON: {e}")
            print(f"Raw response: {generated_text[:500]}...")
            return None

    except requests.exceptions.ConnectionError:
        print("❌ Cannot connect to Ollama. Is it running?")
        print("   Start it with: ollama serve")
        return None
    except requests.exceptions.Timeout:
        print("❌ Request timed out. Try a smaller model or increase timeout.")
        return None
    except Exception as e:
        print(f"❌ Error: {e}")
        return None


def check_ollama_status():
    """Check if Ollama is running and list available models."""
    try:
        response = requests.get('http://localhost:11434/api/tags', timeout=5)
        response.raise_for_status()
        models = response.json()

        print("✅ Ollama is running")
        print(f"📦 Available models: {len(models.get('models', []))}")

        if models.get('models'):
            print("\nInstalled models:")
            for model in models['models']:
                print(f"  - {model['name']}")
        else:
            print("\n⚠️  No models installed. Install one with:")
            print("   ollama pull llama3.2")
            print("   ollama pull mistral")
            print("   ollama pull qwen2.5")

        return True
    except requests.exceptions.ConnectionError:
        print("❌ Ollama is not running")
        print("   Start it with: ollama serve")
        print("   Or install from: https://ollama.ai")
        return False
    except Exception as e:
        print(f"❌ Error checking Ollama: {e}")
        return False


def validate_resume(resume_data: dict, format_name='legacy'):
    """Validate resume data against schema."""
    try:
        from jsonschema import validate, ValidationError
    except ImportError:
        print("⚠️  jsonschema not installed. Skipping validation.")
        print("   Install with: pip install jsonschema")
        return True

    schema = load_schema(format_name)

    try:
        validate(instance=resume_data, schema=schema)
        print("✅ Resume validation passed!")
        return True
    except ValidationError as e:
        print(f"❌ Validation error: {e.message}")
        print(f"   Field path: {'.'.join(str(p) for p in e.path)}")
        return False


def save_resume(resume_data: dict, output_path: str):
    """Save resume as YAML file."""
    with open(output_path, 'w') as f:
        yaml.dump(resume_data, f, default_flow_style=False, sort_keys=False)
    print(f"💾 Saved resume to: {output_path}")


def main():
    """Main demo function."""

    # Example unstructured resume text
    text_resume = """
    Jane Smith
    Senior Software Engineer
    jane.smith@example.com | +1-555-987-6543
    github.com/janesmith | linkedin.com/in/janesmith

    PROFESSIONAL EXPERIENCE

    Tech Innovations Inc | Senior Software Engineer | May 2021 - Present
    - Architected and deployed microservices handling 5M+ requests daily using Go and Kubernetes
    - Reduced infrastructure costs by 45% through optimization and auto-scaling implementation
    - Led team of 8 engineers in successful migration from monolith to microservices
    - Implemented CI/CD pipeline reducing deployment time from 3 hours to 20 minutes

    Digital Solutions Corp | Software Engineer | Jan 2019 - Apr 2021
    - Developed real-time analytics platform processing 100K+ events per second
    - Built RESTful APIs serving 50K+ users with 99.9% uptime
    - Improved database query performance by 70% through indexing and optimization
    - Mentored 3 junior developers in best practices and code review

    StartupXYZ | Full Stack Developer | Jun 2017 - Dec 2018
    - Created responsive web application using React and Node.js
    - Integrated payment processing system handling $2M+ in transactions
    - Implemented automated testing reducing bugs in production by 60%

    EDUCATION

    State University | Bachelor of Science in Computer Science | 2013-2017
    - Graduated with Honors (GPA: 3.85/4.0)
    - Dean's List all semesters
    - Senior thesis on distributed systems

    TECHNICAL SKILLS

    Languages: Python, Go, JavaScript, TypeScript, Java, Rust
    Frameworks & Tools: React, Node.js, Django, Flask, Spring Boot, Kubernetes, Docker
    Cloud & Infrastructure: AWS (EC2, S3, Lambda, ECS, RDS), Terraform, Ansible
    Databases: PostgreSQL, MongoDB, Redis, Elasticsearch
    DevOps: Jenkins, GitHub Actions, ArgoCD, Prometheus, Grafana
    """

    print("=" * 70)
    print("Resume Generator - Ollama LLM Integration Demo")
    print("=" * 70)
    print()

    # Check Ollama status
    if not check_ollama_status():
        return

    print("\n" + "=" * 70)

    # Choose format
    format_name = 'legacy'  # or 'enhanced' or 'json-resume'
    model = 'llama3.2'  # or 'mistral', 'qwen2.5', etc.

    # Generate resume
    resume_data = generate_with_ollama(text_resume, format_name, model)

    if not resume_data:
        print("\n❌ Failed to generate resume")
        print("\nTroubleshooting:")
        print("1. Ensure Ollama is running: ollama serve")
        print(f"2. Pull the model: ollama pull {model}")
        print("3. Try a different model: llama3.2, mistral, qwen2.5")
        return

    print("\n✅ Resume generated successfully!")
    print("\n" + "=" * 70)
    print("Generated Resume Data (YAML):")
    print("=" * 70)
    print(yaml.dump(resume_data, default_flow_style=False, sort_keys=False))
    print("=" * 70)

    # Validate
    print("\n🔍 Validating against schema...")
    is_valid = validate_resume(resume_data, format_name)

    if is_valid:
        # Save to file
        output_path = 'generated-resume.yml'
        save_resume(resume_data, output_path)

        print("\n" + "=" * 70)
        print("Next Steps:")
        print("=" * 70)
        print("📄 Generate PDF:")
        print(f"   ./resume-generator run -i {output_path} -o resume.pdf -t base-latex")
        print("\n🔍 Preview:")
        print(f"   ./resume-generator preview {output_path}")
        print("\n✅ Validate:")
        print(f"   ./resume-generator validate {output_path}")
    else:
        print("\n⚠️  Resume has validation errors. Review and fix before generating PDF.")

    print("\n✨ Demo completed!")
    print("\nOther Ollama models to try:")
    print("  - llama3.2 (fast, good for structured data)")
    print("  - mistral (balanced performance)")
    print("  - qwen2.5 (excellent at following formats)")
    print("  - codellama (optimized for code/structured data)")


if __name__ == '__main__':
    main()
