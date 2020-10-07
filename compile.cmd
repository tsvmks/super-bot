go generate ./app/...

cd app && go build -v -mod=vendor

go mod vendor

go test -mod=vendor ./app/... -coverprofile cover.out

rem .PHONY: lint
rem # lint: runs `golangci-lint`
rem lint:
rem @golangci-lint run ./app/...

cd..

go run -v -mod=vendor app/main.go --dbg "--super=tsvmks"

pause