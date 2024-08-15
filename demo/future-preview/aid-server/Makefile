all:
	make swag
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
	make swag
	make fmt
	go clean -testcache
	# go test -v ./... -cover
	go test ./... -cover

test-function:
	make build
	cd tests && npm run test

fmt:
	go mod tidy
	go fmt ./...

build:
	go mod tidy
	go build -o ./bin/aid ./main.go

build-linux:
	go mod tidy
	CGO_ENABLED=1 GOOS=linux go build -o ./bin/aid ./main.go

docker-build:
	docker buildx build -t aid -f ./Dockerfile --platform linux/amd64 .
	docker tag aid leon1234858/aid-server:latest
	docker push leon1234858/aid-server:latest

deploy:
	echo "create directory ./ssl and put your SSL certificate and key in it..."
	sudo docker-compose up -d

clean:
	rm -rf ./bin
	rm -rf ./data
	rm -rf ./docs