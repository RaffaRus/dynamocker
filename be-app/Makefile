# Golang

dynamocker:
	@ go build -o build/dynamocker cmd/main.go

run: dynamocker
	@ go run cmd/main.go

clear:
	@ rm -rf build/*

test:
# web server tests must be done in series, running them in parallel would lead to concurrency problems during testing
	@ go test ./... -v -p 1 
