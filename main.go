package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	_articleHttpDeliver "gade/srv-gade-point/articles/delivery/http"
	_articleRepo "gade/srv-gade-point/articles/repository"
	_articleUcase "gade/srv-gade-point/articles/usecase"
	"gade/srv-gade-point/middleware"
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
	dbConn := getDBConn()

	defer dbConn.Close()

	dataMigrations(dbConn)

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	ar := _articleRepo.NewPsqlArticleRepository(dbConn)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, timeoutContext)
	_articleHttpDeliver.NewArticlesHandler(e, au)

	e.Start(viper.GetString("server.address"))
}

func getDBConn() *sql.DB {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("postgres://%s%s@%s%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	dbConn, err := sql.Open(`postgres`, connection)

	if err != nil && viper.GetBool(`debug`) {
		fmt.Println(err)
	}

	err = dbConn.Ping()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return dbConn
}

func dataMigrations(dbConn *sql.DB) {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		viper.GetString(`database.user`), driver)

	if err != nil {
		fmt.Println(err)
	}

	migrations.Up()
	defer migrations.Close()
}
