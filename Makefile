build:
	@go build -o bin/api.exe

run: build
	@./bin/api.exe

seed:
	@go run scripts/seed/seed.go

test:
	@clear | go test -v ./... 2>&1 | tee test-output.log