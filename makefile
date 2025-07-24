.PHONY: run test lint migrate

run:
	air


# go install gotest.tools/gotestsum@latest
# gotestsum to see all test cases and create junit for CI/CD 
test:
	gotestsum --format testdox --junitfile junit-report.xml -- \
		-coverpkg="$$(go list ./... | paste -sd "," -)" \
		-covermode=atomic \
		-coverprofile=coverage.out \
		./... && \
	go tool cover -func=coverage.out && \
	go tool cover -html=coverage.out -o coverage.html && \
	open coverage.html

# need to install golangci-lint beforehand
# brew install golangci-lint
lint:
	golangci-lint run --config .golangci.yml

migrate:
	go run migrations/migration.go
