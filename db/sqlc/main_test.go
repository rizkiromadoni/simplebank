package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiromadoni/simplebank/util"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
