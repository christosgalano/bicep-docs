#!/bin/bash

export PATH="/github/home/.azure/bin:${PATH}"
/app/bicep-docs --input "$1" --output "$2" "$3"
