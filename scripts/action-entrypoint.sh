#!/bin/bash

/app/cac "$@" > /tmp/out

# Merge the output into a single line so it can be used as github action output
content=$(awk '{printf "%s\\n", $0}' /tmp/out)
echo "::set-output name=result::${content@Q}"