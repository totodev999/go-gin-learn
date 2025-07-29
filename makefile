.PHONY: run test lint migrate clean build clean_tables

run:
	air


# go install gotest.tools/gotestsum@latest
# gotestsum to see all test cases and create junit for CI/CD 
EXCLUDEDIRS=coverage_exclude_dirs.txt
PKGS=$(shell go list ./... | grep -v -f $(EXCLUDEDIRS) | paste -sd " " -)
COVERPKG=$(shell go list ./... | grep -v -f $(EXCLUDEDIRS) | paste -sd "," -)

test:
	$(MAKE) migrate
	$(MAKE) clean_tables
	@mkdir -p test_coverage
	@echo "-----------test start-----------"
	gotestsum --format testdox --junitfile test_coverage/junit-report.xml -- \
		-coverpkg="$(COVERPKG)" \
		-covermode=atomic \
		-coverprofile=test_coverage/coverage.out \
		./...
	@echo "-----------test done-----------\n\n"
	go tool cover -func=test_coverage/coverage.out | tee test_coverage/totalCoverage.txt
	go tool cover -html=test_coverage/coverage.out -o test_coverage/coverage.html
	open test_coverage/coverage.html
	$(MAKE) clean_tables

# need to install golangci-lint beforehand
# brew install golangci-lint
lint:
	golangci-lint run --config .golangci.yml

migrate:
	@echo "---------------migrate start-----------------"
	go run cmd/migrations/main.go
	@echo "---------------migrate end-----------------\n\n"

clean:
	go clean -cache -testcache

# the command "go build" creates binary for the current environment.(on Linux, binary will be for Linux)
build:
	time go build -o script main.go

clean_tables:
	@echo "---------------clean_tables start-----------------"
	go run cmd/clean_tables/main.go
	@echo "---------------clean_tables end-----------------\n\n"
