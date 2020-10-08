rem Слушаем чат
SET TELEGRAM_GROUP=kristaBIChat

rem id нашего бота
SET TELEGRAM_TOKEN=1332428511:AAGP-c8vs7SdToFLIcW1tmKANReanZCrS2s

go generate ./app/...

cd app && go build -v -mod=vendor

go mod vendor

go test -mod=vendor ./app/... -coverprofile cover.out

rem .PHONY: lint
rem # lint: runs `golangci-lint`
rem lint:
rem @golangci-lint run ./app/...

cd..

go test -mod=vendor ./app/... -coverprofile cover.out

go run -v -mod=vendor app/main.go --dbg "--super=tsvmks"

pause