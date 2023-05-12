SERVER_BINARY_NAME=server_binary
CLIENT_BINARY_NAME=client_binary

build_server:
	GOARCH=amd64 GOOS=darwin go build -o ./build/${SERVER_BINARY_NAME}-darwin cmd/server/main.go
	GOARCH=amd64 GOOS=linux go build -o ./build/${SERVER_BINARY_NAME}-linux cmd/server/main.go
	GOARCH=amd64 GOOS=windows go build -o ./build/${SERVER_BINARY_NAME}-windows cmd/server/main.go

build_client:
	GOARCH=amd64 GOOS=darwin go build -o ./build/${CLIENT_BINARY_NAME}-darwin cmd/client/main.go
	GOARCH=amd64 GOOS=linux go build -o ./build/${CLIENT_BINARY_NAME}-linux cmd/client/main.go
	GOARCH=amd64 GOOS=windows go build -o ./build/${CLIENT_BINARY_NAME}-windows cmd/client/main.go

build_all_simple:
	go build -o ./build/${SERVER_BINARY_NAME} cmd/server/main.go
	go build -o ./build/${CLIENT_BINARY_NAME} cmd/client/main.go

run_server_simple:
	go run ./build/${SERVER_BINARY_NAME}

run_client_simple:
	go run ./build/${CLIENT_BINARY_NAME}

test:
	go test ./...

clean:
	go clean
	rm ./build/${SERVER_BINARY_NAME}-darwin ./build/${SERVER_BINARY_NAME}-linux ./build/${SERVER_BINARY_NAME}-windows
	rm ./build/${CLIENT_BINARY_NAME}-darwin ./build/${CLIENT_BINARY_NAME}-linux ./build/${CLIENT_BINARY_NAME}-windows
	rm ./build/${SERVER_BINARY_NAME} ./build/${CLIENT_BINARY_NAME}
