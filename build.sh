#!/bin/bash -x
# VERSION=$1
TEMP=$(mktemp -d build-XXXXX)
trap 'rm -rf $TEMP' EXIT

VERSION=$(git describe --tags --always --dirty)
flags=(-ldflags="-X github.com/bitte-ein-bit/songbeamer-helper/cmd.version=$VERSION -X github.com/bitte-ein-bit/songbeamer-helper/cmd.updateURL=https://software.ec-pfuhl.de/")

export GOARCH=amd64
for GOOS in darwin windows; do
    echo "Building $VERSION for $GOOS"
    go build "${flags[@]}" -o "$TEMP/$GOOS-$GOARCH" main.go
done

unset GOARCH GOOS
echo "Making self update"
go-selfupdate "$TEMP" "$VERSION"
aws s3 sync public/ s3://software.ec-pfuhl.de/songbeamer-helper/ --delete --acl public-read
