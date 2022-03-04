package database

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/logger"
	"os"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func GetDBConn() (*sql.DB, *DbBun) {
	var dbBun DbBun
	dbHost := os.Getenv(`DB_HOST`)
	dbPort := os.Getenv(`DB_PORT`)
	dbUser := os.Getenv(`DB_USER`)
	dbPass := os.Getenv(`DB_PASS`)
	dbName := os.Getenv(`DB_NAME`)

	connection := fmt.Sprintf("postgres://%s%s@%s%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := sql.Open(`postgres`, connection)

	if err != nil {
		logger.Make(nil, nil).Debug(err)

		return nil, nil
	}

	err = dbConn.Ping()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		os.Exit(1)
	}

	// bun connection init
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connection)))
	dbBun.DB = bun.NewDB(sqldb, pgdialect.New())

	if os.Getenv(`DB_LOGGER`) == "true" {
		dbBun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return dbConn, &dbBun
}
