#!/bin/bash

export PATH="/usr/bin/az:${PATH}"
/app/bicep-docs --input "$1" --output "$2" "$3"
