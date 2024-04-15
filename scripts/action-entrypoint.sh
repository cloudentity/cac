#!/bin/sh

/app/cac "$@" > /tmp/out

if [ $? -ne 0 ]; then
    # save status to variable
    status=$?
fi

# Merge the output into a single line so it can be used as github action output
content=$(awk '{printf "%s\\n", $0}' /tmp/out)
echo "::set-output name=result::$content"

if [ $status -ne 0 ]; then
    exit $status
fi