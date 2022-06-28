# Init variables
GOBIN := $(shell go env GOBIN)

# Keep this at the top so that it is default when `make` is called.
# This is used by Travis CI.
coverage.txt: gen-mocks
	go test -race -coverprofile=coverage.txt.tmp -covermode=atomic ./...
	# Remove coverage for mock files
	cat coverage.txt.tmp | grep -v "/mock_" > coverage.txt
	rm coverage.txt.tmp

view-cover: clean coverage.txt
	go tool cover -html=coverage.txt
test: build gen-mocks
	go test ./...
build:
	go build ./...
install: build
	go install ./...
inspect: build $(GOBIN)/golangci-lint
	golangci-lint run --skip-files /mock_
update:
	go get -u ./...
gen-mocks: $(GOBIN)/mockery
	mockery --all --case underscore --inpackage
pre-commit: update clean coverage.txt inspect
	go mod tidy
clean:
	rm -f coverage.txt $(GOBIN)/sportrank

# Needed tools
$(GOBIN)/golangci-lint:
	$(MAKE) install-tools
$(GOBIN)/mockery:
	$(MAKE) install-tools
install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.42.1
	rm -rf ./v1.42.1
	go install github.com/vektra/mockery/v2@latest