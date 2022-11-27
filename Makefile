GOTEST ?= go test
GOTOOL ?= go tool

setup:
	go install github.com/rakyll/gotest

.PHONY: test
test:
	go test -v ./...

.PHONY: fuzz
fuzz:
	go test -fuzztime=10s -fuzz .

COVERAGE_OUT = coverage.out
$(COVERAGE_OUT): *.go
	$(GOTEST) -cover -coverprofile=coverage.out -v ./...

coverage-report: $(COVERAGE_OUT)
	$(GOTOOL) cover -html=coverage.out