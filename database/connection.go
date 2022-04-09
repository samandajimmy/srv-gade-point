package database

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"os"
	"reflect"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	Logger   string

	Sql *sql.DB
	Bun *DbBun
}

func NewDbConn(dbConfig DbConfig) *DbConfig {
	// check if dbConfig empty or not
	if reflect.DeepEqual(dbConfig, DbConfig{}) {
		dbConfig = DbConfig{
			Host:     os.Getenv(`DB_HOST`),
			Port:     os.Getenv(`DB_PORT`),
			Username: os.Getenv(`DB_USER`),
			Password: os.Getenv(`DB_PASS`),
			Name:     os.Getenv(`DB_NAME`),
			Logger:   os.Getenv(`DB_LOGGER`),
		}
	}

	var dbBun DbBun
	var err error

	connection := fmt.Sprintf("postgres://%s%s@%s%s/%s?sslmode=disable",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)

	dbConfig.Sql, err = sql.Open(`postgres`, connection)

	if err != nil {
		logger.Make(nil).Fatal(err)
	}

	err = dbConfig.Sql.Ping()

	if err != nil {
		logger.Make(nil).Fatal(err)
	}

	// bun connection init
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connection)))
	dbBun.DB = bun.NewDB(sqldb, pgdialect.New())

	if dbConfig.Logger == "true" {
		dbBun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	dbConfig.Bun = &dbBun
	return &dbConfig
}

func (dbc *DbConfig) Migrate(dbConn *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})

	if err != nil {
		logger.Make(nil).Fatal(err)
	}

	migrationPath := "migrations"

	if os.Getenv(`APP_PATH`) != "" {
		migrationPath = os.Getenv(`APP_PATH`) + "/" + migrationPath
	}

	// check migration dir existence
	_, err = os.Stat(migrationPath)

	if os.IsNotExist(err) {
		migrationPath = "migrations"
	}

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		dbc.Username, driver)

	if err != nil {
		logger.Make(nil).Fatal(err)
	}

	if err := migrations.Up(); err != nil && err.Error() != models.ErrMigrateNoChange.Error() {
		logger.Make(nil).Fatal(err)
	}

	return migrations
}
