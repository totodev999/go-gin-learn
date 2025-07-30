.PHONY: run test lint migrate clean build clean_tables test_race

run:
	air


# go install gotest.tools/gotestsum@latest
# gotestsum to see all test cases and create junit for CI/CD 
EXCLUDEDIRS=coverage_exclude_dirs.txt
PKGS=$(shell go list ./... | grep -v -f $(EXCLUDEDIRS) | paste -sd " " -)
COVERPKG=$(shell go list ./... | grep -v -f $(EXCLUDEDIRS) | paste -sd "," -)

test:
	$(MAKE) clean_tables
	$(MAKE) migrate
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

# It's not determined whether race condition can be detected or not, it depends on timing.
# And cache can cause test to be skipped, so clean cache beforehand.
test_race:
	$(MAKE) clean_tables
	$(MAKE) migrate
	$(MAKE) clean
	@echo "Running tests with race detector..."
	@out="$$(gotestsum --format testdox -- -race ./... 2>&1)"; \
	echo "$$out"; \
	if echo "$$out" | grep -q "DATA RACE"; then \
		echo "\033[0;31m❗️Race condition detected!\033[0m"; \
	else \
		echo "\033[0;32m✅ No race condition found.\033[0m"; \
	fi

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
