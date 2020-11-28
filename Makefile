.PHONY: test

test:
	go test ./pkg/...

build:
	go build cmd/ww.go