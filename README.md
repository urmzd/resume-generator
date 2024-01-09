# Generate Resumes 

Generate elegant $\LaTeX$ resumes effortlessly from a single source using a single command.

## Usage

To seamlessly interact with this tool, it is *highly recommended* that users utilize the Docker entrypoint and leverage volumes appropriately.

```bash
# Step 1. Build the Docker image.
docker build -t generate-resumes .
# Step 2. Run the Docker container with two bind mounts.
# -v "$(pwd)/outputs" ensures access to all generated resumes.
# -v "$(pwd)/inputs" allows adding custom templates without rebuilding. Place your custom config files here.
docker run -v "$(pwd)/outputs:/outputs" -v "$(pwd)/inputs:/inputs" generate-resumes inputs/example.toml -o /outputs
```

:warning: If you prefer to avoid Docker, consult the Dockerfile in the project's root to set up your local environment. However, exercise caution, as maintaining and building the required dependencies may pose challenges.