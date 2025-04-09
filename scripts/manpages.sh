#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run ./cmd/surgeon man | gzip -c -9 >manpages/surgeon.1.gz
