#!/bin/bash

# Function to extract flag
extract_flag() {
  if [[ "$1" == *=false ]]; then
    echo ""
  elif [[ "$1" == *=true ]]; then
    echo "${1%=true}"
  else
    return 1
  fi
}

# GitHub Actions related
verbose=$(extract_flag "$3")
return_code=$?
if [[ $return_code -eq 1 ]]; then
  echo "Error: Invalid argument for --verbose (true | false)"
  exit 1
fi

/app/bicep-docs --input "$1" --output "$2" "$verbose"
