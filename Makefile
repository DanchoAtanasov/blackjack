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

build:
	docker-compose build
