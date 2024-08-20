APP_NAME=aid-server

all:
	make swag
	go run main.go

swag:
	swag fmt
	swag init

cli-tools:
	go get -u github.com/swaggo/echo-swagger@latest
	go install github.com/swaggo/swag/cmd/swag@latest

build:
	go build -o bin/$(APP_NAME) main.go

clean:
	rm -rf bin/$(APP_NAME)
	rm -rf data/ldb/