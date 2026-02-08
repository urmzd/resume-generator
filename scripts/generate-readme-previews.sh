#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
input_file="${root_dir}/assets/example_resumes/software_engineer.yml"
output_dir="${root_dir}/assets/example_results"
binary="${root_dir}/resume-generator"
templates=("modern-html" "modern-latex")

if ! command -v pdftoppm >/dev/null 2>&1; then
  echo "Error: pdftoppm (poppler-utils) is required for PDF-to-PNG conversion." >&2
  exit 1
fi

echo "Building CLI binary..."
(cd "$root_dir" && go build -trimpath -ldflags="-s -w" -o "$binary" .)

mkdir -p "$output_dir"

for template in "${templates[@]}"; do
  tmp_dir="$(mktemp -d)"
  echo "Generating PDF for ${template}..."
  "$binary" run -i "$input_file" -o "$tmp_dir" -t "$template"

  pdf_path="$(find "$tmp_dir" -type f -name "*.pdf" -print -quit)"
  if [ -z "$pdf_path" ]; then
    rm -rf "$tmp_dir"
    echo "Error: failed to locate PDF output for ${template}." >&2
    exit 1
  fi

  target_path="${output_dir}/${template}.png"
  echo "Converting ${pdf_path} -> ${target_path}"
  pdftoppm -png -singlefile -r 300 "$pdf_path" "${target_path%.png}"

  rm -rf "$tmp_dir"
done

echo "Preview images updated in ${output_dir}"
