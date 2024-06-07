go mod tidy
gofumpt -l -w .
golangci-lint run
