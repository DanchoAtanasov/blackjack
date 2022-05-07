build:
	go build -o /usr/local/bin/app

run:
	./app

dev:
	go run main.go

test:
	go test -v blackjack/models

coverage:
	go test --cover blackjack/models