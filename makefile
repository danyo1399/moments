install:
	go install gotest.tools/gotestsum@latest
test:
	gotestsum --packages="./..."
prepare:
	go mod tidy
	go test ./...
testwatch:
	gotestsum --watch --packages="./..."
