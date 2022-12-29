GOTEST ?= go test
GOTOOL ?= go tool

setup:
	go get github.com/rakyll/gotest
	go install github.com/rakyll/gotest
	go get golang.org/x/tools/cmd/godoc
	go install golang.org/x/tools/cmd/godoc
	go mod tidy

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

doc: setup
	godoc -http=:6060
