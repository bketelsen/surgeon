#!/bin/sh
set -e

SED="sed"
if which gsed >/dev/null 2>&1; then
	SED="gsed"
fi
NEXT=`svu n`
wholething="# surgeon <small>$NEXT</small>"
# update this directory to the default value
# of the `--output` flag on the doc generation command
# and at the end of the script too 
rm -rf www/docs/cmd/*.md
go run ./cmd/surgeon gendocs
"$SED" \
	-i'' \
	-e 's/SEE ALSO/See also/g' \
	-e 's/^## /# /g' \
	-e 's/^### /## /g' \
	-e 's/^#### /### /g' \
	-e 's/^##### /#### /g' \
	./docs/surgeon*.md
echo $NEXT
"$SED" \
	-i'' \
	 "/v[0-9]\+\.[0-9]\+\.[0-9]/c $wholething" \
	./docs/_coverpage.md
