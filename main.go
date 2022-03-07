package main

import (
	"bufio"
	"database/sql"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"net/http"
	"os"
	"time"

	_campaignHttpDelivery "gade/srv-gade-point/campaigns/delivery/http"
	_campaignTrxHttpDelivery "gade/srv-gade-point/campaigntrxs/delivery/http"
	_pHistoryHttpDelivery "gade/srv-gade-point/pointhistories/delivery/http"
	_referralsHttpDelivery "gade/srv-gade-point/referrals/delivery/http"
	_referralTrxHttpDelivery "gade/srv-gade-point/referraltrxs/delivery/http"
	_rewardHttpDelivery "gade/srv-gade-point/rewards/delivery/http"
	_rewardTrxHttpDelivery "gade/srv-gade-point/rewardtrxs/delivery/http"
	_tokenHttpDelivery "gade/srv-gade-point/tokens/delivery/http"
	_userHttpDelivery "gade/srv-gade-point/users/delivery/http"
	_voucherCodeHttpDelivery "gade/srv-gade-point/vouchercodes/delivery/http"
	_voucherHttpDelivery "gade/srv-gade-point/vouchers/delivery/http"

	"github.com/carlescere/scheduler"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var ech *echo.Echo

func init() {
	// helper.NewHelper()
	ech = echo.New()
	ech.Debug = true
	config.LoadEnv()
	logger.Init()
}

func main() {
	dbConn, dbBun := database.GetDBConn()
	migrate := dataMigrations(dbConn)

	defer dbConn.Close()
	defer migrate.Close()

	echoGroup := models.EchoGroup{
		Admin: ech.Group("/admin"),
		API:   ech.Group("/api"),
		Token: ech.Group("/token"),
	}

	// load all middlewares
	middleware.InitMiddleware(ech, echoGroup)

	// PING
	ech.GET("/ping", ping)

	// Load all repositories
	repos := config.NewRepositories(dbConn, dbBun)
	// Load all usecases
	usecases := config.NewUsecases(repos)

	// load all handlers
	_tokenHttpDelivery.NewTokensHandler(echoGroup, usecases.TokenUseCase)
	_pHistoryHttpDelivery.NewPointHistoriesHandler(echoGroup, usecases.PHistoryUseCase)
	_rewardTrxHttpDelivery.NewRewardTrxHandler(echoGroup, usecases.RewardTrxUC)
	_referralTrxHttpDelivery.NewReferralTrxHandler(echoGroup, usecases.ReferralTrxUseCase)
	_voucherCodeHttpDelivery.NewVoucherCodesHandler(echoGroup, usecases.VoucherCodeUseCase)
	_voucherHttpDelivery.NewVouchersHandler(echoGroup, usecases.VoucherUseCase)
	_rewardHttpDelivery.NewRewardHandler(echoGroup, usecases.RewardUseCase)
	_campaignHttpDelivery.NewCampaignsHandler(echoGroup, usecases.CampaignUseCase)
	_campaignTrxHttpDelivery.NewCampaignTrxsHandler(echoGroup, usecases.CampaignTrxUseCase,
		usecases.CampaignUseCase)
	_userHttpDelivery.NewUserHandler(echoGroup, usecases.UserUseCase)
	_referralsHttpDelivery.NewReferralsHandler(echoGroup, usecases.ReferralsUseCase)

	// Run every day.
	updateStatusBasedOnStartDate(usecases.CampaignUseCase, usecases.VoucherUseCase)

	// run refresh reward trx
	usecases.RewardUseCase.RefreshTrx()

	_ = ech.Start(":" + os.Getenv(`PORT`))

}

func updateStatusBasedOnStartDate(cmp campaigns.UseCase, vcr vouchers.VUsecase) {
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

	response := models.Response{}
	response.Status = models.StatusSuccess
	response.Message = "PONG!!"
	response.Data = text

	return echTx.JSON(http.StatusOK, response)
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
