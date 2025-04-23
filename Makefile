DB_URL=postgresql://root:secret@localhost:5433/personal_financial?sslmode=disable

createdb:
	docker exec -it postgres createdb --username=root --owner=root personal_financial

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root personal_financial

migrateup: 
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -databate "$(DB_URL)" -verbose down

sqlc:
	docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run .

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/sangketkit01/simple-bank/db/sqlc Store
.PHONY: