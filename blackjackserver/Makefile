build:
	go build -o /usr/local/bin/app
	cd audit_client && go build -o /usr/local/bin/auditclient

run:
	./app

dev:
	go run main.go

aud:
	echo $(SEED) $(SESSION_ID)
	cd audit_client && go build -o ac . && ./ac
	

test:
	go test -v ./...

coverage:
	go test --cover ./...

clean:
	rm -rf logs/*
	rm -rf results/*
	rm -rf audit/*
