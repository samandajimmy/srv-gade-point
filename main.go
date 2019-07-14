package main

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	_campaignHttpDelivery "gade/srv-gade-point/campaigns/delivery/http"
	_campaignRepository "gade/srv-gade-point/campaigns/repository"
	_campaignUseCase "gade/srv-gade-point/campaigns/usecase"
	_campaignTrxHttpDelivery "gade/srv-gade-point/campaigntrxs/delivery/http"
	_campaignTrxRepository "gade/srv-gade-point/campaigntrxs/repository"
	_campaignTrxUseCase "gade/srv-gade-point/campaigntrxs/usecase"
	_pHistoryHttpDelivery "gade/srv-gade-point/pointhistories/delivery/http"
	_pHistoryRepository "gade/srv-gade-point/pointhistories/repository"
	_pHistoryUseCase "gade/srv-gade-point/pointhistories/usecase"
	_quotaRepository "gade/srv-gade-point/quotas/repository"
	_quotaUseCase "gade/srv-gade-point/quotas/usecase"
	_rewardHttpDelivery "gade/srv-gade-point/rewards/delivery/http"
	_rewardRepository "gade/srv-gade-point/rewards/repository"
	_rewardUseCase "gade/srv-gade-point/rewards/usecase"
	_tagRepository "gade/srv-gade-point/tags/repository"
	_tagUseCase "gade/srv-gade-point/tags/usecase"
	_tokenHttpDelivery "gade/srv-gade-point/tokens/delivery/http"
	_tokenRepository "gade/srv-gade-point/tokens/repository"
	_tokenUseCase "gade/srv-gade-point/tokens/usecase"
	_userHttpDelivery "gade/srv-gade-point/users/delivery/http"
	_userRepository "gade/srv-gade-point/users/repository"
	_userUseCase "gade/srv-gade-point/users/usecase"
	_voucherCodeHttpDelivery "gade/srv-gade-point/vouchercodes/delivery/http"
	_voucherCodeRepository "gade/srv-gade-point/vouchercodes/repository"
	_voucherCodeUseCase "gade/srv-gade-point/vouchercodes/usecase"
	_voucherHttpDelivery "gade/srv-gade-point/vouchers/delivery/http"
	_voucherRepository "gade/srv-gade-point/vouchers/repository"
	_voucherUseCase "gade/srv-gade-point/vouchers/usecase"

	"github.com/carlescere/scheduler"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var ech *echo.Echo

func init() {
	ech = echo.New()
	ech.Debug = true
	loadEnv()
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			tmp := strings.Split(f.File, "/")
			filename := tmp[len(tmp)-1]
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

	logrus.SetFormatter(formatter)
	logrus.SetLevel(logrus.DebugLevel)
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
		Admin: ech.Group("/admin"),
		API:   ech.Group("/api"),
		Token: ech.Group("/token"),
	}

	// load all middlewares
	middleware.InitMiddleware(ech, echoGroup)

	// PING
	ech.GET("/ping", ping)

	// TOKEN
	tokenRepository := _tokenRepository.NewPsqlTokenRepository(dbConn)
	tokenUseCase := _tokenUseCase.NewTokenUseCase(tokenRepository, timeoutContext)
	_tokenHttpDelivery.NewTokensHandler(echoGroup, tokenUseCase)

	// TAG
	tagRepository := _tagRepository.NewPsqlTagRepository(dbConn)
	tagUseCase := _tagUseCase.NewTagUseCase(tagRepository)

	// QUOTA
	quotaRepository := _quotaRepository.NewPsqlQuotaRepository(dbConn)
	quotaUseCase := _quotaUseCase.NewQuotaUseCase(quotaRepository)

	// REWARD
	rewardRepository := _rewardRepository.NewPsqlRewardRepository(dbConn)
	campaignRepository := _campaignRepository.NewPsqlCampaignRepository(dbConn, rewardRepository)
	rewardUseCase := _rewardUseCase.NewRewardUseCase(rewardRepository, campaignRepository, tagUseCase, quotaUseCase)
	_rewardHttpDelivery.NewRewardHandler(echoGroup, rewardUseCase)

	// CAMPAIGN
	campaignUseCase := _campaignUseCase.NewCampaignUseCase(campaignRepository, rewardUseCase)
	_campaignHttpDelivery.NewCampaignsHandler(echoGroup, campaignUseCase)

	// CAMPAIGNTRX
	campaignTrxRepository := _campaignTrxRepository.NewPsqlCampaignTrxRepository(dbConn)
	campaignTrxUseCase := _campaignTrxUseCase.NewCampaignTrxUseCase(campaignTrxRepository)
	_campaignTrxHttpDelivery.NewCampaignTrxsHandler(echoGroup, campaignTrxUseCase, campaignUseCase)

	// POINTHISTORY
	pHistoryRepository := _pHistoryRepository.NewPsqlPointHistoryRepository(dbConn)
	pHistoryUseCase := _pHistoryUseCase.NewPointHistoryUseCase(pHistoryRepository)
	_pHistoryHttpDelivery.NewPointHistoriesHandler(echoGroup, pHistoryUseCase)

	// VOUCHER
	voucherRepository := _voucherRepository.NewPsqlVoucherRepository(dbConn)
	voucherUseCase := _voucherUseCase.NewVoucherUseCase(voucherRepository, campaignRepository, pHistoryRepository)
	_voucherHttpDelivery.NewVouchersHandler(echoGroup, voucherUseCase)

	// VOUCHERCODE
	voucherCodeRepository := _voucherCodeRepository.NewPsqlVoucherCodeRepository(dbConn)
	voucherCodeUseCase := _voucherCodeUseCase.NewVoucherCodeUseCase(voucherCodeRepository, voucherRepository)
	_voucherCodeHttpDelivery.NewVoucherCodesHandler(echoGroup, voucherCodeUseCase)

	// USER
	userRepository := _userRepository.NewPsqlUserRepository(dbConn)
	userUseCase := _userUseCase.NewUserUseCase(userRepository, timeoutContext)
	_userHttpDelivery.NewUserHandler(echoGroup, userUseCase)

	// Run every day.
	updateStatusBasedOnStartDate(campaignUseCase, voucherUseCase)

	ech.Start(":" + os.Getenv(`PORT`))

}

func updateStatusBasedOnStartDate(cmp campaigns.UseCase, vcr vouchers.UseCase) {
	scheduler.Every().Day().At(os.Getenv(`STATUS_UPDATE_TIME`)).Run(func() {
		t := time.Now()
		logrus.Debug("Run Scheduler! @", t)

		// CAMPAIGN
		cmp.UpdateStatusBasedOnStartDate()

		// VOUCHER
		vcr.UpdateStatusBasedOnStartDate()
	})
}

func ping(echTx echo.Context) error {
	res := echTx.Response()
	rid := res.Header().Get(echo.HeaderXRequestID)
	params := map[string]interface{}{"rid": rid}

	requestLogger := logrus.WithFields(logrus.Fields{"params": params})

	requestLogger.Info("Start to ping server.")
	response := models.Response{}
	response.Status = models.StatusSuccess
	response.Message = "PONG!!"

	requestLogger.Info("End of ping server.")

	return echTx.JSON(http.StatusOK, response)
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
		logrus.Debug(err)
	}

	err = dbConn.Ping()

	if err != nil {
		logrus.Fatal(err)
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
		logrus.Debug(err)
	}

	if err := migrations.Up(); err != nil {
		logrus.Debug(err)
	}

	return migrations
}

func loadEnv() {
	// check .env file existence
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		return
	}

	err := godotenv.Load()

	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	return
}
