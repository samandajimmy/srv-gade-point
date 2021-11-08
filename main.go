package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"net/http"
	"os"
	"strconv"
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
	_referralTrxHttpDelivery "gade/srv-gade-point/referraltrxs/delivery/http"
	_referralTrxRepository "gade/srv-gade-point/referraltrxs/repository"
	_referralTrxUseCase "gade/srv-gade-point/referraltrxs/usecase"
	_rewardHttpDelivery "gade/srv-gade-point/rewards/delivery/http"
	_rewardRepository "gade/srv-gade-point/rewards/repository"
	_rewardUseCase "gade/srv-gade-point/rewards/usecase"
	_rewardTrxHttpDelivery "gade/srv-gade-point/rewardtrxs/delivery/http"
	_rewardTrxRepository "gade/srv-gade-point/rewardtrxs/repository"
	_rewardTrxUC "gade/srv-gade-point/rewardtrxs/usecase"
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
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var ech *echo.Echo

func init() {
	ech = echo.New()
	ech.Debug = true
	loadEnv()
	logger.Init()
}

func main() {
	dbConn, dbBun := getDBConn()
	migrate := dataMigrations(dbConn)

	defer dbConn.Close()
	defer migrate.Close()

	contextTimeout, err := strconv.Atoi(os.Getenv(`CONTEXT_TIMEOUT`))

	if err != nil {
		logger.Make(nil, nil).Debug(err)
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

	// POINTHISTORY
	pHistoryRepository := _pHistoryRepository.NewPsqlPointHistoryRepository(dbConn)
	pHistoryUseCase := _pHistoryUseCase.NewPointHistoryUseCase(pHistoryRepository)
	_pHistoryHttpDelivery.NewPointHistoriesHandler(echoGroup, pHistoryUseCase)

	// REWARDTRX
	rewardTrxRepository := _rewardTrxRepository.NewPsqlRewardTrxRepository(dbConn)
	rewardTrxUseCase := _rewardTrxUC.NewRewardtrxUseCase(rewardTrxRepository)
	_rewardTrxHttpDelivery.NewRewardTrxHandler(echoGroup, rewardTrxUseCase)

	// REFERRALTRX
	referralTrxRepository := _referralTrxRepository.NewPsqlReferralTrxRepository(dbConn)
	referralTrxUseCase := _referralTrxUseCase.NewReferralTrxUseCase(referralTrxRepository)
	_referralTrxHttpDelivery.NewReferralTrxHandler(echoGroup, referralTrxUseCase)

	// QUOTA
	quotaRepository := _quotaRepository.NewPsqlQuotaRepository(dbConn)
	quotaUseCase := _quotaUseCase.NewQuotaUseCase(quotaRepository, rewardTrxUseCase)

	// VOUCHER
	voucherRepository := _voucherRepository.NewPsqlVoucherRepository(dbConn, dbBun)

	// VOUCHERCODE
	voucherCodeRepository := _voucherCodeRepository.NewPsqlVoucherCodeRepository(dbConn)
	voucherCodeUseCase := _voucherCodeUseCase.NewVoucherCodeUseCase(voucherCodeRepository, voucherRepository)
	_voucherCodeHttpDelivery.NewVoucherCodesHandler(echoGroup, voucherCodeUseCase)

	// REWARD
	rewardRepository := _rewardRepository.NewPsqlRewardRepository(dbConn, quotaRepository, tagRepository)
	campaignRepository := _campaignRepository.NewPsqlCampaignRepository(dbConn, rewardRepository)
	voucherUseCase := _voucherUseCase.NewVoucherUseCase(voucherRepository, campaignRepository, pHistoryRepository, voucherCodeRepository)
	_voucherHttpDelivery.NewVouchersHandler(echoGroup, voucherUseCase)
	rewardUseCase := _rewardUseCase.NewRewardUseCase(rewardRepository, campaignRepository, tagUseCase, quotaUseCase, voucherUseCase, voucherCodeRepository, rewardTrxRepository, referralTrxRepository)
	_rewardHttpDelivery.NewRewardHandler(echoGroup, rewardUseCase)

	// CAMPAIGN
	campaignUseCase := _campaignUseCase.NewCampaignUseCase(campaignRepository, rewardUseCase)
	_campaignHttpDelivery.NewCampaignsHandler(echoGroup, campaignUseCase)

	// CAMPAIGNTRX
	campaignTrxRepository := _campaignTrxRepository.NewPsqlCampaignTrxRepository(dbConn)
	campaignTrxUseCase := _campaignTrxUseCase.NewCampaignTrxUseCase(campaignTrxRepository)
	_campaignTrxHttpDelivery.NewCampaignTrxsHandler(echoGroup, campaignTrxUseCase, campaignUseCase)

	// USER
	userRepository := _userRepository.NewPsqlUserRepository(dbConn)
	userUseCase := _userUseCase.NewUserUseCase(userRepository, timeoutContext)
	_userHttpDelivery.NewUserHandler(echoGroup, userUseCase)

	// Run every day.
	updateStatusBasedOnStartDate(campaignUseCase, voucherUseCase)

	// run refresh reward trx
	rewardUseCase.RefreshTrx()

	_ = ech.Start(":" + os.Getenv(`PORT`))

}

func updateStatusBasedOnStartDate(cmp campaigns.UseCase, vcr vouchers.UseCase) {
	_, _ = scheduler.Every().Day().At(os.Getenv(`STATUS_UPDATE_TIME`)).Run(func() {
		t := time.Now()
		logger.Make(nil, nil).Debug("Run Scheduler! @", t)

		// CAMPAIGN
		_ = cmp.UpdateStatusBasedOnStartDate()

		// VOUCHER
		_ = vcr.UpdateStatusBasedOnStartDate()
	})
}

func ping(echTx echo.Context) error {
	res := echTx.Response()
	rid := res.Header().Get(echo.HeaderXRequestID)
	params := map[string]interface{}{"rid": rid}
	file, err := os.Open("latest_commit_hash")

	if err != nil {
		logger.Make(nil, nil).Debug("failed to open")

	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	logger.Make(echTx, params).Info("Start to ping server.")
	response := models.Response{}
	response.Status = models.StatusSuccess
	response.Message = "PONG!!"
	response.Data = text
	logger.Make(echTx, params).Info("End of ping server.")

	return echTx.JSON(http.StatusOK, response)
}

func getDBConn() (*sql.DB, *bun.DB) {
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
	}

	err = dbConn.Ping()

	if err != nil {
		logger.Make(nil, nil).Debug(err)
		os.Exit(1)
	}

	// bun connection init
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connection)))
	dbBun := bun.NewDB(sqldb, pgdialect.New())

	if os.Getenv(`DB_LOGGER`) == "true" {
		dbBun.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return dbConn, dbBun
}

func dataMigrations(dbConn *sql.DB) *migrate.Migrate {
	driver, _ := postgres.WithInstance(dbConn, &postgres.Config{})

	migrations, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		os.Getenv(`DB_USER`), driver)

	if err != nil {
		logger.Make(nil, nil).Debug(err)
	}

	if err := migrations.Up(); err != nil {
		logger.Make(nil, nil).Debug(err)
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
		logger.Make(nil, nil).Fatal("Error loading .env file")
	}
}
