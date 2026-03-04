#!/bin/bash

VERSION=$(git describe --tags --always --dirty)

# Read credentials from external file
CREDS=$(cat .googlecloud.json | base64 | tr -d '\n')

# Read Church Tools credentials - these should be set in environment or loaded from a secure file
# For now, we'll use empty defaults if not set
CT_USERID="${CT_USERID:-}"
CT_TOKEN="${CT_TOKEN:-}"

flags=(-ldflags="-X github.com/bitte-ein-bit/songbeamer-helper/cmd.version=$VERSION -X github.com/bitte-ein-bit/songbeamer-helper/cmd.updateURL=https://software.ec-pfuhl.de/ -X github.com/bitte-ein-bit/songbeamer-helper/log.credsJSON=$CREDS -X github.com/bitte-ein-bit/songbeamer-helper/churchtools.userid=$CT_USERID -X github.com/bitte-ein-bit/songbeamer-helper/churchtools.token=$CT_TOKEN")

go run "${flags[@]}" main.go "$@"