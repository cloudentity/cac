#!/bin/sh

/app/cac "$@" > /tmp/out

content=$(awk '{printf "%s\\n", $0}' /tmp/out)
echo "::set-output name=result::$content"