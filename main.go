package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbConn, err := sql.Open(`postgres`, getDBConn())

	if err != nil && viper.GetBool(`debug`) {
		fmt.Println(err)
	}

	err = dbConn.Ping()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"minijarvis", driver)

	if err != nil {
		fmt.Println(err)
	}

	m.Steps(2)
	defer m.Close()
	defer dbConn.Close()

	e := echo.New()
	e.Start(viper.GetString("server.address"))
}

func getDBConn() string {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("postgres://%s%s@%s%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)
	return connection
}
