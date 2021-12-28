build:
	@npm run build

serve:
	@go run .\server\server.go

dev:
	@make build
	@make serve

release:
	@make build
	@if not exist .\dist-ssr mkdir .\dist-ssr
	@xcopy .\dist .\dist-ssr\dist\ /E /I /Q /Y

release-windows:
	@make release
	@cd .\server && go build -o ..\dist-ssr\server.exe

release-linux:
	@make release
	@set CGO_ENABLED=0
	@set GOOS=linux
	@set GOARCH=amd64
	@cd .\server && go build -o ..\dist-ssr\server.o