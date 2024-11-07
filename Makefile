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

proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --experimental_allow_proto3_optional --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simplebank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

redis:
	docker run -d --name redis -p 6379:6379 redis:7.4-alpine

test:
	go test -v -cover ./...

start:
	go run main.go

.PHONY: migrateup migratedown sqlc mock dbschema dbdocs proto redis test start