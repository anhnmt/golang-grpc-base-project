
docker.build:
	docker compose -f docker-compose.yml build

docker.push:
	docker compose -f docker-compose.yml push

docker.dev: docker.build docker.push

docker.pull:
	docker compose -f docker-compose.yml pull

docker.run:
	docker compose -f docker-compose.yml up -d --force-recreate

docker.up: docker.pull docker.run

docker.down:
	docker compose -f docker-compose.yml down

