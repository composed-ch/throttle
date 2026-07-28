.PHONY: cover format test ver

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html
	rm -f cover.out

# setup: go install golang.org/x/tools/cmd/goimports@latest
format:
	goimports -w .

test:
	go test ./...

vet:
	go vet ./...