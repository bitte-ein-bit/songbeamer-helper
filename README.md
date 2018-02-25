# Building

     go get
     go build -ldflags="-X github.com/bitte-ein-bit/songbeamer-helper/cmd.version=0.5 -X github.com/bitte-ein-bit/songbeamer-helper/cmd.updateURL=http://some-update-url" -o songbeamer-helper.exe main.go
     go-selfupdate songbeamer-helper.exe 0.5
     aws s3 sync public/ s3://update-path/songbeamer-helper/ --delete --acl public-read
