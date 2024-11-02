migrateup:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simplebank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/rizkiromadoni/simplebank/db/sqlc Store

dbschema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

dbdocs:
	dbdocs build doc/db.dbml

test:
	go test -v -cover ./...

start:
	go run main.go

.PHONY: migrateup migratedown sqlc mock dbschema dbdocs test start