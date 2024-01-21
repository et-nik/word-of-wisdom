lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run --timeout 5m0s ./...

build-all: build-server build-client

build-server:
	docker build -t knik/wofw-server -f ./build/server.Dockerfile .

build-client:
	docker build -t knik/wofw-client -f ./build/client.Dockerfile .

run-server:
	docker network create wofw-net 2> /dev/null | true

	docker rm -f wofw-server 2> /dev/null | true

	docker run --rm \
		--name wofw-server \
		--net wofw-net \
		knik/wofw-server:latest

run-client:
	docker network create wofw-net 2> /dev/null | true

	docker rm -f wofw-client 2> /dev/null | true

	docker run --rm \
		--name wofw-client \
		--net wofw-net \
		knik/wofw-client:latest \
		--host=wofw-server --port=9100

clean:
	docker rm -f wofw-server wofw-client | true
	docker network remove wofw-net | true