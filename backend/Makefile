.PHONY: build-macos bump

build-macos:
	GOOS=darwin GOARCH=arm64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-darwin-arm64.dylib main.go
	GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-darwin-amd64.dylib main.go

bump:
	go get -u
	go mod tidy
