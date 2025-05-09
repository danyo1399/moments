install:
	go install gotest.tools/gotestsum@latest
test:
	go test ./...
prepare:
	go mod tidy
	go test ./...
testwatch:
	gotestsum --watch --packages="./..."
