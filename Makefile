TAG=demoservice

all: test compile docker

test:
	go test -v -race -cover ./...
	go test -bench=. -run=Benchmark ./...

compile:
	CGO_ENABLED=0 go build -o bin/service .

run: run_db
	go run .

docker:
	docker build -t $(TAG) .

clean:
	docker rmi -f $(TAG)  > /dev/null 2>&1 || true
	rm -f bin/service


DB_UP:=$(shell docker ps | grep 'demoservice-db')

run_db:
ifndef DB_UP
	docker run --name demoservice-db -d \
	-e POSTGRES_PASSWORD=postgres \
	-e POSTGRES_DB=demoservice \
	-p 5432:5432 \
	-v $(PWD)/sql:/docker-entrypoint-initdb.d \
	postgres
	sleep 5
endif

stop_db:
	docker rm -f demoservice-db
