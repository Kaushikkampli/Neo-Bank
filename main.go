package main

import (
	"database/sql"
	"log"

	"github.com/kaushikkampli/neobank/api"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/utils"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot instantiate server", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
