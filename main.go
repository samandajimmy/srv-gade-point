package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConn := getDBConn()
	dataMigrations(dbConn)

	defer dbConn.Close()

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)

	contextTimeout, err := strconv.Atoi(os.Getenv(`CONTEXT TIME`))

	if err != nil {
		fmt.Println(err)
	}
	timeoutContext := time.Duration(contextTimeout) * time.Second

	ar := _articleRepo.NewPsqlArticleRepository(dbConn)
	au := _articleUcase.NewArticleUsecase(ar, timeoutContext)
	_articleHttpDeliver.NewArticlesHandler(e, au)

	e.Start(os.Getenv(`SERVER_PORT`))
}

func getDBConn() *sql.DB {
	dbHost := os.Getenv(`DB_HOST`)
	dbPort := os.Getenv(`DB_PORT`)
	dbUser := os.Getenv(`DB_USER`)
	dbPass := os.Getenv(`DB_PASS`)
	dbName := os.Getenv(`DB_NAME`)

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
		os.Getenv(`DB_USER`), driver)

	if err != nil {
		fmt.Println(err)
	}

	migrations.Up()
	defer migrations.Close()
}
