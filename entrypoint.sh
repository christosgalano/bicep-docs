#!/bin/bash

ln -s /github/home/.azure/bin/az /usr/bin/az

/app/bicep-docs --input "$1" --output "$2" "$3"
