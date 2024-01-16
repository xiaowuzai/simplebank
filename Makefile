postgres:
	docker run --network bank_network --name postgres16 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=123456 -d postgres:16-alpine

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

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go	

mock:
	mockgen -package mockdb  -destination db/mock/store.go github.com/xiaowuzai/simplebank/db/sqlc Store

docker:
	docker build -f Dockerfile -t simplebank:latest .

.PHONY: postgres createdb dropdb migrateup migratedown  migrateup1 migratedown1 sqlc test server mock docker 