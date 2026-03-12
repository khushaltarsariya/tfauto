#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
templates_dir="${repo_root}/templates"

if [[ ! -d "${templates_dir}" ]]; then
  echo "templates directory not found: ${templates_dir}" >&2
  exit 1
fi

work_dir="$(mktemp -d)"
cleanup() {
  rm -rf "${work_dir}"
}
trap cleanup EXIT

echo "Validating templates from ${templates_dir}"

for template_path in "${templates_dir}"/*; do
  [[ -d "${template_path}" ]] || continue

  template_name="$(basename "${template_path}")"
  temp_template_dir="${work_dir}/${template_name}"

  echo
  echo "==> ${template_name}"

  cp -R "${template_path}" "${temp_template_dir}"

  terraform fmt -check -recursive "${temp_template_dir}"

  pushd "${temp_template_dir}" >/dev/null
  terraform init -backend=false -input=false
  terraform validate
  popd >/dev/null
done

echo
echo "All templates validated successfully."
