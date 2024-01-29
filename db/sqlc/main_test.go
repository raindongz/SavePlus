package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/randongz/save_plus/utils"
)


var testQueries *Queries
var conn *pgxpool.Pool

func TestMain(m *testing.M){
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err = pgxpool.New(context.Background(), config.DBSource )
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}