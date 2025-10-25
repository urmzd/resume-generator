# Makefile for resume-generator project

# Variables
ORGANIZATION   := urmzd
VERSION        ?= latest

# Docker image names
BASE_IMAGE_NAME := resume-generator-base
BASE_IMAGE_TAG  := $(ORGANIZATION)/$(BASE_IMAGE_NAME):$(VERSION)
APP_IMAGE_TAG   := $(ORGANIZATION)/resume-generator:$(VERSION)

# Directories
OUTPUTS_DIR   := outputs
INPUTS_DIR    := inputs
EXAMPLES_DIR  := assets/example_inputs
TEMPLATES_DIR := templates
TEMPLATE_NAME ?= modern-html

# ────────────────────────────────────────────────────────────────────────────────
# Build the base Docker image (TeX tooling)
build-base:
	@echo "Building base Docker image $(BASE_IMAGE_TAG)..."
	docker build \
	  --file base.Dockerfile \
	  --tag $(BASE_IMAGE_TAG) \
	  .

# Push the base image
push-base:
	@echo "Pushing base image to Docker Hub..."
	docker push $(BASE_IMAGE_TAG)

# Build the application image (multi-stage pulls builder + base)
build:
	@echo "Building application Docker image $(APP_IMAGE_TAG)..."
	docker build \
	  --tag $(APP_IMAGE_TAG) \
	  .

# Initialize directories & copy examples
init:
	@echo "Initializing $(INPUTS_DIR) and $(OUTPUTS_DIR)..."
	mkdir -p $(INPUTS_DIR) $(OUTPUTS_DIR)
	cp -r $(EXAMPLES_DIR)/* $(INPUTS_DIR)/

# Run the resume-generator inside Docker
run:
	@echo "Running resume-generator..."
	$(eval KEEPTEX_FLAG := $(if $(KEEPTEX),-k,))
	docker run --rm \
	  -v "$(CURDIR)/$(INPUTS_DIR):/inputs" \
	  -v "$(CURDIR)/$(OUTPUTS_DIR):/outputs" \
	  -v "$(CURDIR)/$(TEMPLATES_DIR):/templates" \
	  $(APP_IMAGE_TAG) run \
	  -i /inputs/$(FILENAME) \
	  -o /outputs/$(basename $(FILENAME)).pdf \
	  $(KEEPTEX_FLAG) \
	  -t $(TEMPLATE_NAME)

# Exec an arbitrary command in the image
exec:
	@echo "Executing in resume-generator container..."
	docker run --rm -it \
	  -v "$(CURDIR)/$(INPUTS_DIR):/inputs" \
	  -v "$(CURDIR)/$(TEMPLATES_DIR):/templates" \
	  $(APP_IMAGE_TAG) \
	  $(CMD)

# Start an interactive shell
shell:
	@echo "Launching shell in resume-generator container..."
	docker run --rm -it \
	  -v "$(CURDIR)/$(INPUTS_DIR):/inputs" \
	  -v "$(CURDIR)/$(OUTPUTS_DIR):/outputs" \
	  -v "$(CURDIR)/$(TEMPLATES_DIR):/templates" \
	  --entrypoint /bin/bash \
	  $(APP_IMAGE_TAG)

# Clean generated outputs & inputs
clean:
	@echo "Cleaning up $(INPUTS_DIR) and $(OUTPUTS_DIR)..."
	rm -rf $(INPUTS_DIR) $(OUTPUTS_DIR)

# Run all examples through the pipeline
examples: init
	@echo "Running all example resumes..."
	@for f in $(EXAMPLES_DIR)/*; do \
	  filename=$$(basename $$f); \
	  make run FILENAME=$$filename KEEPTEX=true; \
	done
