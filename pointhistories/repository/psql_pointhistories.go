package repository

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/pointhistories"

	"github.com/labstack/echo"
)

type psqlPointHistoryRepository struct {
	Conn *sql.DB
}

// NewPsqlPointHistoryRepository will create an object that represent the pointhistories.Repository interface
func NewPsqlPointHistoryRepository(Conn *sql.DB) pointhistories.Repository {
	return &psqlPointHistoryRepository{Conn}
}

func (psqlRepo *psqlPointHistoryRepository) CountUsers(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	queryCounter := `SELECT COUNT(*) counter FROM
					( SELECT CIF
						FROM point_histories
						group by CIF
					) tmp;`

	err := psqlRepo.Conn.QueryRow(queryCounter).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (psqlRepo *psqlPointHistoryRepository) GetUsers(c echo.Context, payload map[string]interface{}) ([]models.PointHistory, error) {
	var pntHstrs []models.PointHistory
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	paging := ""

	query := `SELECT ph.CIF, sum(tmp.conv_point_amt) point_amount FROM point_histories ph left join (select id, 
		CASE WHEN transaction_type = 'K' THEN 0-point_amount ELSE point_amount end conv_point_amt from point_histories) tmp on tmp.id = ph.id
		group by CIF order by point_amount desc`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += paging + ";"
	rows, err := psqlRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pntHstry models.PointHistory

		err = rows.Scan(
			&pntHstry.CIF,
			&pntHstry.PointAmount,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		pntHstrs = append(pntHstrs, pntHstry)
	}

	return pntHstrs, nil
}

func (psqlRepo *psqlPointHistoryRepository) CountUserPointHistory(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	where := ""
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `select COUNT(*) counter from point_histories where CIF = $1`

	if payload["startDateRg"].(string) != "" {
		where += " and transaction_date::timestamp::date >= '" + payload["startDateRg"].(string) + "'"
	}

	if payload["endDateRg"].(string) != "" {
		where += " and transaction_date::timestamp::date <= '" + payload["endDateRg"].(string) + "'"
	}

	query += where + ";"
	err := psqlRepo.Conn.QueryRow(query, payload["CIF"].(string)).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (psqlRepo *psqlPointHistoryRepository) GetUserPointHistory(c echo.Context, payload map[string]interface{}) ([]models.PointHistory, error) {
	var dataHistory []models.PointHistory
	where := ""
	paging := ""
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	startDateRg := payload["startDateRg"].(string)
	endDateRg := payload["endDateRg"].(string)

	// TODO update the query to a reward base
	query := `select
				ph.id,
				ph.CIF,
				ph.point_amount,
				ph.transaction_type,
				ph.transaction_date,
				coalesce(ph.reff_core, '') reff_core,
				coalesce(ph.campaign_id, 0) campaign_id,
				coalesce(c.name, '') campaign_name,
				coalesce(c.description, '') campaign_description,
				coalesce(ph.voucher_code_id, 0) voucher_code_id,
				coalesce(pc.promo_code, '') promo_code,
				coalesce(pc.voucher_id, 0) voucher_id,
				coalesce(v.name, '') voucher_name,
				coalesce(v.description, '') voucher_description
			from point_histories ph
			left join campaigns c on ph.campaign_id = c.id
			left join voucher_codes pc on pc.id = ph.voucher_code_id
			left join vouchers v on pc.voucher_id = v.id
			where ph.CIF = $1`

	if startDateRg != "" {
		where += " and ph.transaction_date::timestamp::date >= '" + startDateRg + "'"
	}

	if endDateRg != "" {
		where += " and ph.transaction_date::timestamp::date <= '" + endDateRg + "'"
	}

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += where + " order by ph.transaction_date desc" + paging + ";"
	rows, err := psqlRepo.Conn.Query(query, payload["CIF"].(string))

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var ph models.PointHistory
		var campaign models.Campaign
		var voucherCodes models.VoucherCode
		var voucher models.Voucher

		err = rows.Scan(
			&ph.ID,
			&ph.CIF,
			&ph.PointAmount,
			&ph.TransactionType,
			&ph.TransactionDate,
			&ph.RefCore,
			&campaign.ID,
			&campaign.Name,
			&campaign.Description,
			&voucherCodes.ID,
			&voucherCodes.PromoCode,
			&voucher.ID,
			&voucher.Name,
			&voucher.Description,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		// TODO change this to a reward base
		// if campaign.ID != 0 {
		// ph.Reward = &campaign
		// }

		if voucherCodes.ID != 0 {
			ph.VoucherCode = &voucherCodes
			ph.VoucherCode.ID = 0 // remove promo codes ID from the response
		}

		if voucher.ID != 0 {
			ph.VoucherCode.Voucher = &voucher
		}

		dataHistory = append(dataHistory, ph)
	}

	return dataHistory, nil
}

func (psqlRepo *psqlPointHistoryRepository) GetUserPoint(c echo.Context, CIF string) (float64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var pointDebet float64
	var pointKredit float64
	queryDebet := `SELECT coalesce(sum(point_amount), 0) as debet FROM point_histories
		WHERE cif = $1 AND transaction_type = 'D' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`
	err := psqlRepo.Conn.QueryRow(queryDebet, CIF).Scan(&pointDebet)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	queryKredit := `SELECT coalesce(sum(point_amount), 0) as debet FROM point_histories
		WHERE cif = $1 AND transaction_type = 'K' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`
	err = psqlRepo.Conn.QueryRow(queryKredit, CIF).Scan(&pointKredit)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	pointAmount := pointDebet - pointKredit

	return pointAmount, nil
}
