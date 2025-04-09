#!/usr/bin/env bash
FILES=$(git diff --cached --name-only --diff-filter=ACMR)

gofmt -l -w .
golangci-lint run --new --fix

git add $FILES
