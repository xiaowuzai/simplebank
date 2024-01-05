package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/xiaowuzai/simplebank/api"
	db "github.com/xiaowuzai/simplebank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:123456@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
