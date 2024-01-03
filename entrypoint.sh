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
result=""
summary=$(extract_flag "$8")
return_code=$?
if [[ $return_code -eq 1 ]]; then
  echo "Error: Invalid argument for --summary (true | false)"
  exit 1
fi

# Get command, path and include-preview
command="$1"
path="$2"

include_preview=$(extract_flag "$3")
return_code=$?
if [[ $return_code -eq 1 ]]; then
  echo "Error: Invalid argument for --include-preview (true | false)"
  exit 1
fi

# Get the appropriate arguments for the command
if [[ "$command" == "scan" ]]; then
    outdated=$(extract_flag "$4")
    return_code=$?
    if [[ $return_code -eq 1 ]]; then
      echo "Error: Invalid argument for --outdated (true | false)"
      exit 1
    fi
    output="$5"
    if [[ "$output" != *=normal && "$output" != *=table && "$output" != *=markdown ]]; then
      echo "Error: Invalid argument for --output (normal | table | markdown)"
      exit 1
    fi

    result=$(eval "/app/bicep-docs $command $path $include_preview $outdated $output")

elif [[ "$command" == "update" ]]; then
    in_place=$(extract_flag "$6")
    return_code=$?
    if [[ $return_code -eq 1 ]]; then
      echo "Error: Invalid argument for --in-place (true | false)"
      exit 1
    fi
    silent=$(extract_flag "$7")
    return_code=$?
    if [[ $return_code -eq 1 ]]; then
      echo "Error: Invalid argument for --silent (true | false)"
      exit 1
    fi

    echo "In place: $in_place"
    echo "Silent: $silent"

    result=$(eval "/app/bicep-docs $command $path $include_preview $in_place $silent")
else
    echo "Error: Command not found (scan/update)"
    exit 1
fi

echo "::debug::\$result: $result"

EOF=$(dd if=/dev/urandom bs=15 count=1 status=none | base64)
echo "result<<$EOF" >> "$GITHUB_OUTPUT"
echo "$result" >> "$GITHUB_OUTPUT"
echo "$EOF" >> "$GITHUB_OUTPUT"

if [[ "$summary" == "--summary" ]]; then
  if [[ "$command" == "scan" && "$output" != *=markdown ]]; then
    result=$(eval "/app/bicep-docs $command $path $include_preview $outdated --output=markdown")
  fi
  if [[ "$command" == "update" ]]; then
    echo "## Update results" >> "$GITHUB_STEP_SUMMARY"
  else
    echo "## Scan results" >> "$GITHUB_STEP_SUMMARY"
  fi
  echo "$result" >> "$GITHUB_STEP_SUMMARY"
  echo "---" >> "$GITHUB_STEP_SUMMARY"
fi
