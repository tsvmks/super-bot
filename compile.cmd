go generate ./app/...

cd app && go build -v -mod=vendor

go mod vendor

go test -mod=vendor ./app/... -coverprofile cover.out

#.PHONY: lint
## lint: runs `golangci-lint`
#lint:
# @golangci-lint run ./app/...

go run -v -mod=vendor app/main.go --dbg "--super=tsvmks"

pause