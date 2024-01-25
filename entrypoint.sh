#!/bin/bash

az config set bicep.use_binary_from_path=False

/app/bicep-docs --input "$1" --output "$2" "$3"
