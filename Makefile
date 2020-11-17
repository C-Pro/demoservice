TAG=demoservice

all: test compile docker

test:
	go test -v -race -cover ./...
	go test -bench=. -run=Benchmark ./...

compile:
	CGO_ENABLED=0 go build -o bin/service .

run:
	go run .

docker:
	docker build -t $(TAG) .

clean:
	docker rmi -f $(TAG)  > /dev/null 2>&1 || true
	rm -f bin/service
