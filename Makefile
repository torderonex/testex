build:
	docker-compose build testex

run:
	docker-compose up testex

test:
	go test -v ./...
