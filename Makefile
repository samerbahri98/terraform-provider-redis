include .env
export

.PHONY: test
test:
	GOPROXY="off" GOFLAGS="-mod=vendor" go test -count=1 -v ./...
	GOPROXY="off" GOFLAGS="-mod=vendor" go vet ./...
