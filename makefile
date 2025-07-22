.PHONY: test-cover

test-cover:
	go test ./... -coverpkg="$$(go list ./... | paste -sd "," -)" -covermode=atomic -coverprofile=coverage.out && \
	go tool cover -func=coverage.out | grep total && \
	go tool cover -html=coverage.out