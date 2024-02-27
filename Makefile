postgres-network:
	docker run --network bank_network --name postgres16 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=123456 -d postgres:16-alpine

postgres:
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=123456 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=admin --owner=admin simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://admin:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://admin:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://admin:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://admin:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go	

local:
	ENVIRONMENT=local go run main.go	

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/xiaowuzai/simplebank/db/sqlc Store

docker:
	docker build -f Dockerfile -t simplebank:latest .

proto: 
	rm -rf pb/*
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger  --openapiv2_opt=allow_merge=true \
	--openapiv2_opt=merge_file_name=simple_bank --openapiv2_opt=json_names_for_fields=false \
    proto/*.proto

evans:
	evans --host localhost --port 9090 --reflection rep

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres-network postgres createdb dropdb migrateup migratedown  migrateup1 migratedown1 sqlc test local server \
	mock docker db_docs db_schema proto redis