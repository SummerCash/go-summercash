#!/bin/bash
# Note: this script was taken from the go-spacemesh repository.
# All credits for this script, validate_lint.sh, go to the spacemeshos development team.

pkgs=`go list ./...`

output=`golint $pkgs 2>&1 | grep -vE "_mock|_test"`
if [ $(echo -n "$output" |  wc -l) -ne 0 ]; then
    echo -n "$output";
    exit 1;
fi