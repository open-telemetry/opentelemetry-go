#!/bin/bash

set -euo pipefail

top_dir='.'
if [[ $# -gt 0 ]]; then
    top_dir="${1}"
fi

p=$(pwd)
mod_dirs=()
mapfile -t mod_dirs < <(find "${top_dir}" -type f -name 'go.mod' -exec dirname {} \; | sort)

for mod_dir in "${mod_dirs[@]}"; do
    cd "${mod_dir}"
    main_dirs=()
    mapfile -t main_dirs < <(go list --find -f '{{.Name}}|{{.Dir}}' ./... | grep '^main|' | cut -f 2- -d '|')
    for main_dir in "${main_dirs[@]}"; do
        echo ".${main_dir#${p}}"
    done
    cd "${p}"
done
