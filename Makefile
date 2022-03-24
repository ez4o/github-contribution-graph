build:
	@npm run build

serve:
	@go run .\main.go

dev:
	@make build
	@make serve

release:
	@make build
	@if not exist .\dist-ssr\dist mkdir .\dist-ssr\dist
	@copy .\config.json .\dist-ssr\config.json
	@copy .\dist\index.html .\dist-ssr\dist\index.html

release-windows:
	@make release
	@go build -o .\dist-ssr\server.exe .

release-linux:
	@make release
	@REM avoid white space after env variable
	@set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&& go build -o .\dist-ssr\server.o .