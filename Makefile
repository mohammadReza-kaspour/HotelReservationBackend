build:
	@go build -o bin/api.exe

run: build
	@./bin/api.exe

test:
	@clear | go test -v ./... 2>&1 | tee test-output.log