build:
	go build -o /usr/local/bin/app

run:
	./app

dev:
	go run main.go

test:
	go test -v ./...

coverage:
	go test --cover ./...
