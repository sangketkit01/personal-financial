DB_URL=postgresql://root:secret@localhost:5433/personal_financial?sslmode=disable

createdb:
	docker exec -it postgres createdb --username=root --owner=root personal_financial

dropdb:
	docker exec -it postgres dropdb --username=root --owner=root personal_financial

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup: 
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run .

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/sangketkit01/simple-bank/db/sqlc Store


.PHONY: createdb dropdb new_migration migrateup migratedown sqlc test server mock