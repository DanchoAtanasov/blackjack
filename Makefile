up:
	docker-compose up

rup:
	docker-compose build && docker-compose up

down:
	docker-compose down

stop:
	docker-compose stop

db:
	docker exec -it database psql -U postgres

blackjack:
	docker exec -it blackjackserver /bin/sh

redis:
	docker exec -it redis redis-cli

build:
	docker-compose build

aud:
	echo "Running audit"
	docker-compose up -d
	docker exec -it blackjackserver ./auditclient
	docker-compose stop

env: 
	docker-compose run blackjackserver env
