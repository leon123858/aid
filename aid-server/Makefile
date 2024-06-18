all:
	make fmt
	make run

init:
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/swaggo/echo-swagger

swag:
	swag fmt
	swag init -g ./server/router.go

run:
	go run ./main.go server

test:
	go clean -testcache
	#go test -v ./... -cover
	go test ./... -cover
fmt:
	go mod tidy
	go fmt ./...

build:
	go mod tidy
	go build -o ./bin/aid ./main.go

clean:
	rm -rf ./bin