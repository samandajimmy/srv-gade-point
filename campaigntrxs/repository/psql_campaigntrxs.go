package repository

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/campaigntrxs"
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type psqlCampaignTrxRepository struct {
	Conn *sql.DB
}

// NewPsqlCampaignTrxRepository will create an object that represent the campaigntrxs.Repository interface
func NewPsqlCampaignTrxRepository(Conn *sql.DB) campaigntrxs.Repository {
	return &psqlCampaignTrxRepository{Conn}
}

func (psqlRepo *psqlCampaignTrxRepository) CountUsers(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	queryCounter := `SELECT COUNT(*) counter FROM
					( SELECT user_id
						FROM campaign_transactions
						group by user_id
					) tmp;`

	err := psqlRepo.Conn.QueryRow(queryCounter).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (psqlRepo *psqlCampaignTrxRepository) GetUsers(c echo.Context, payload map[string]interface{}) ([]models.CampaignTrx, error) {
	var cmpTrxs []models.CampaignTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	paging := ""

	query := `SELECT ct.user_id, sum(tmp.conv_point_amt) point_amount FROM campaign_transactions ct left join (select id, 
		CASE WHEN transaction_type = 'K' THEN 0-point_amount ELSE point_amount end conv_point_amt from campaign_transactions) tmp on tmp.id = ct.id
		group by user_id order by point_amount desc`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += paging + ";"
	rows, err := psqlRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	for rows.Next() {
		var cmpTrx models.CampaignTrx

		err = rows.Scan(
			&cmpTrx.CIF,
			&cmpTrx.PointAmount,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		cmpTrxs = append(cmpTrxs, cmpTrx)
	}

	return cmpTrxs, nil
}
