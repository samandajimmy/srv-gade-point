package main

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/models"
	"log"
	"os"
	"strconv"
	"time"

	_campaignHttpDelivery "gade/srv-gade-point/campaigns/delivery/http"
	_campaignRepository "gade/srv-gade-point/campaigns/repository"
	_campaignUseCase "gade/srv-gade-point/campaigns/usecase"
	_tokenHttpDelivery "gade/srv-gade-point/tokens/delivery/http"
	_voucherHttpDelivery "gade/srv-gade-point/vouchers/delivery/http"
	_voucherRepository "gade/srv-gade-point/vouchers/repository"
	_voucherUseCase "gade/srv-gade-point/vouchers/usecase"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var ech *echo.Echo

func init() {
	ech = echo.New()

	isDebug, _ := strconv.ParseBool(os.Getenv(`DEBUG_APP`))

	if isDebug {
		loadEnv()
	}

	ech = echo.New()

	// setup PUBLIC DIRECTORY
	ech.Static(os.Getenv(`VOUCHER_PATH`), os.Getenv(`VOUCHER_ROUTE_PATH`))

}

func main() {
	dbConn := getDBConn()
	migrate := dataMigrations(dbConn)

	defer dbConn.Close()
	defer migrate.Close()

	contextTimeout, err := strconv.Atoi(os.Getenv(`CONTEXT_TIMEOUT`))

	if err != nil {
		fmt.Println(err)
	}

	timeoutContext := time.Duration(contextTimeout) * time.Second

	echoGroup := models.EchoGroup{
		ech.Group("/admin"),
		ech.Group("/api"),
		ech.Group("/token"),
	}

	// load all middlewares
	middleware.InitMiddleware(ech, echoGroup)

	// TOKEN
	_tokenHttpDelivery.NewTokensHandler(echoGroup)

	// CAMPAIGN
	campaignRepository := _campaignRepository.NewPsqlCampaignRepository(dbConn)
	campaignUseCase := _campaignUseCase.NewCampaignUseCase(campaignRepository, timeoutContext)
	_campaignHttpDelivery.NewCampaignsHandler(echoGroup, campaignUseCase)

	// VOUCHER
	voucherRepository := _voucherRepository.NewPsqlVoucherRepository(dbConn)
	voucherUseCase := _voucherUseCase.NewVoucherUseCase(voucherRepository, campaignRepository, timeoutContext)
	_voucherHttpDelivery.NewVouchersHandler(ech, voucherUseCase)

	ech.Start(":" + os.Getenv(`PORT`))
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

	if err != nil {
		fmt.Println(err)
	}

	err = dbConn.Ping()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return dbConn
}

func dataMigrations(dbConn *sql.DB) *migrate.Migrate {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		os.Getenv(`DB_USER`), driver)

	if err != nil {
		fmt.Println(err)
	}

	if err := migrations.Up(); err != nil {
		fmt.Println(err)
	}

	return migrations
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
