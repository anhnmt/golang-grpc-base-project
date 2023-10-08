
docker.dev:
	docker compose -f docker-compose.yml build --push

docker.up:
	docker compose -f docker-compose.yml up -d --force-recreate --pull

docker.down:
	docker compose -f docker-compose.yml down

wire.gen:
	wire ./...

ent.new:
	ent new User

ent.gen:
	ent generate --feature privacy,entql,schema/snapshot,sql/lock,sql/upsert ./ent/schema

go.install:
	go install github.com/google/wire/cmd/wire@v0.5.0
	go install entgo.io/ent/cmd/ent@v0.12.4