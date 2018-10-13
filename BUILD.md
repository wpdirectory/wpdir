# Build Steps

## Generate Protobuf Files

Example for the `/search/` package:
```protoc -I=. -I=%GOPATH%/src -I=%GOPATH%/src/github.com/gogo/protobuf/protobuf --gogoslick_out=. search.proto```

Generated files are added to git, so this only needs to be done if updating the .proto files.

## Build React Project

From the `/web/` folder:
```npm run build```

## Embed the React files into Go

From root `go generate`, this will embed the files from `web/build/` into the Go files.

## Build the Final Binary

This is a windows PowerShell command, for other OS you will need to adjust it yourself.
`go install -tags=prod -ldflags="-s -w -X main.version=1.1.0 -X main.commit=$(git rev-parse --verify HEAD) -X main.date=$((Get-Date).toString("yyyy-MM-dd"))"`