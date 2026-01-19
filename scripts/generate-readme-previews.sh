#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
input_file="/app/assets/example_resumes/software_engineer.yml"
output_dir="${root_dir}/assets/example_results"
image_tag="resume-generator:readme-previews"
templates=("modern-html" "modern-latex")

if ! command -v docker >/dev/null 2>&1; then
  echo "Error: docker is required to generate README previews." >&2
  exit 1
fi

if ! command -v magick >/dev/null 2>&1 && ! command -v convert >/dev/null 2>&1 && ! command -v pdftoppm >/dev/null 2>&1; then
  echo "Error: install ImageMagick (magick/convert) or poppler-utils (pdftoppm) for PDF -> PNG conversion." >&2
  exit 1
fi

mkdir -p "$output_dir"

echo "Building preview image..."
docker build -t "$image_tag" "$root_dir" >/dev/null

for template in "${templates[@]}"; do
  container_name="resume-preview-${template//[^a-zA-Z0-9]/-}"
  echo "Generating PDF for ${template}..."
  docker create --name "$container_name" "$image_tag" \
    run -i "$input_file" -o /outputs -t "$template" >/dev/null
  docker start -a "$container_name" >/dev/null

  tmp_dir="$(mktemp -d)"
  if ! docker cp "${container_name}:/outputs" "$tmp_dir" >/dev/null 2>&1; then
    docker rm -f "$container_name" >/dev/null 2>&1 || true
    echo "Error: failed to copy outputs from container for ${template}." >&2
    exit 1
  fi

  docker rm -f "$container_name" >/dev/null 2>&1 || true

  pdf_path="$(find "$tmp_dir/outputs" -type f -name "*_resume.pdf" -print -quit)"
  if [ -z "$pdf_path" ]; then
    rm -rf "$tmp_dir"
    echo "Error: failed to locate PDF output for ${template}." >&2
    exit 1
  fi

  target_path="${output_dir}/${template}.png"
  echo "Converting ${pdf_path} -> ${target_path}"
  if command -v magick >/dev/null 2>&1; then
    magick -density 300 "${pdf_path}[0]" -quality 90 "$target_path"
  elif command -v convert >/dev/null 2>&1; then
    convert -density 300 "${pdf_path}[0]" -quality 90 "$target_path"
  else
    pdftoppm -png -singlefile -r 300 "$pdf_path" "${target_path%.png}"
  fi

  rm -rf "$tmp_dir"
done

echo "Preview images updated in ${output_dir}"
