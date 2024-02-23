package apis

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/randongz/save_plus/db/sqlc"
	"github.com/randongz/save_plus/utils"
	"github.com/stretchr/testify/require"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func NewTestServer(t *testing.T) *Server {
	config, err := utils.LoadConfig("../.")
	if err != nil {
		log.Fatal("loading config error", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("connect to db failed", err)
	}

	server, err := NewServer(config, *db.NewStore(conn))
	require.NoError(t, err)

	return server
}
