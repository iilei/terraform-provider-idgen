.PHONY: update-deps tidy test commit

# Update all dependencies
update-deps:
	go get -u ./...

# Tidy up the module
tidy:
	go mod tidy

# Run tests
test:
	go test ./...

# Full process: update, tidy, test
update: update-deps tidy test
