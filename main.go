package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kaushikkampli/neobank/api"
	db "github.com/kaushikkampli/neobank/db/sqlc"
	"github.com/kaushikkampli/neobank/utils"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	config, err := utils.LoadConfig(".")

	fmt.Println("here")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
