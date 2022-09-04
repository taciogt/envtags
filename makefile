GOTEST ?= go test
GOTOOL ?= go tool

.PHONY: test
test:
	go test -v ./...

coverage:
	$(GOTEST) -cover -coverprofile=coverage.out -v ./...
	$(GOTOOL) cover -html=coverage.out