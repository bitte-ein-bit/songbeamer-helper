#!/bin/bash -x
# VERSION=$1
TEMP=$(mktemp -d build-XXXXX)
trap 'rm -rf $TEMP' EXIT

VERSION=$(git describe --tags --always --dirty)

# Read credentials from external file
CREDS=$(cat .googlecloud.json | base64 | tr -d '\n')

# Read Church Tools credentials - these should be set in environment or loaded from a secure file
# For now, we'll use empty defaults if not set
CT_USERID="${CT_USERID:-}"
CT_TOKEN="${CT_TOKEN:-}"

flags=(-ldflags="-X github.com/bitte-ein-bit/songbeamer-helper/cmd.version=$VERSION -X github.com/bitte-ein-bit/songbeamer-helper/cmd.updateURL=https://software.ec-pfuhl.de/ -X github.com/bitte-ein-bit/songbeamer-helper/log.credsJSON=$CREDS -X github.com/bitte-ein-bit/songbeamer-helper/churchtools.userid=$CT_USERID -X github.com/bitte-ein-bit/songbeamer-helper/churchtools.token=$CT_TOKEN")

export GOARCH=amd64
for GOOS in darwin windows; do
    export GOOS
    echo "Building $VERSION for $GOOS"
    go build "${flags[@]}" -o "$TEMP/$GOOS-$GOARCH" main.go
    file "$TEMP/$GOOS-$GOARCH"
done

ls -lah "$TEMP"
md5sum "$TEMP"/*
sleep 3
unset GOARCH GOOS
echo "Making self update"
go-selfupdate "$TEMP" "$VERSION"
export AWS_PROFILE=privat
find public -name '*darwin*' -delete
aws s3 sync public/ s3://software.ec-pfuhl.de/songbeamer-helper/ --delete
