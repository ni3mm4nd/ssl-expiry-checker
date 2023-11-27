GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/ssl-expiry-checker-amd64-windows.exe
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/ssl-expiry-checker-amd64-linux
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/ssl-expiry-checker-arm64-macos