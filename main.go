package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/randongz/save_plus/apis"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/utils"
)

func main() {

	config, err := utils.LoadConfig(".") // current folder. main and app.env are in the same folder
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := apis.NewServer(config, *store)
	if err != nil {
		log.Fatal("canot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
