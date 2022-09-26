.PHONY: build serve dev release release-windows release-linux

build:
	pnpm run build

serve:
	go run .\main.go

dev: build serve

release: build
	if not exist .\dist-ssr\dist mkdir .\dist-ssr\dist
	copy .\config.json .\dist-ssr\config.json
	copy .\dist\index.html .\dist-ssr\dist\index.html

release-windows: release
	go build -o .\dist-ssr\server.exe .

release-linux: release
	set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&& go build -o .\dist-ssr\server.o .