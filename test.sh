#!/bin/bash

VERSION=$(git describe --tags --always --dirty)

# Read credentials from external file
CREDS=$(cat .googlecloud.json | base64 | tr -d '\n')

flags=(-ldflags="-X github.com/bitte-ein-bit/songbeamer-helper/cmd.version=$VERSION -X github.com/bitte-ein-bit/songbeamer-helper/cmd.updateURL=https://software.ec-pfuhl.de/ -X github.com/bitte-ein-bit/songbeamer-helper/log.credsJSON=$CREDS")

go run "${flags[@]}" main.go "$@"