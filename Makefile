DB_USER=root
DB_PASSWORD=secret
DB_HOST=localhost
DB_PORT=5432
DB_NAME=simple_bank

postgres-network:
	docker run --network bank_network --name postgres16 -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:16-alpine

postgres:
	docker run --name postgres16 -p $(DB_PORT):$(DB_PORT) \
	-e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD)  -e POSTGRES_DB=$(DB_NAME) \
	-d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it postgres16 dropdb $(DB_NAME)

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc: 
	sqlc generate

test:
	go test -v -short -cover ./...

server:
	go run main.go	

local:
	ENVIRONMENT=local go run main.go	

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/xiaowuzai/simplebank/db/sqlc Store
	mockgen -package mockworker  -destination worker/mock/distributor.go github.com/xiaowuzai/simplebank/worker TaskDistributor

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
	--openapiv2_opt=merge_file_name=$(DB_NAME) --openapiv2_opt=json_names_for_fields=false \
    proto/*.proto

evans:
	evans --host $(DB_HOST) --port 9090 --reflection rep

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres-network postgres createdb dropdb new_migration migrateup migratedown  migrateup1 migratedown1 sqlc test local server \
	mock docker db_docs db_schema proto redis