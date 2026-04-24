.PHONY: build-macos build-linux build-windows bump

build-macos:
	GOOS=darwin GOARCH=arm64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-darwin-arm64.dylib .
	GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-darwin-amd64.dylib .

build-linux:
	GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-linux-amd64.so .

build-windows:
	GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o ./whatsmeow/whatsmeow-windows-amd64.dll .

bump:
	go get -u
	go mod tidy
