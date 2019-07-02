package repository

import (
	"database/sql"
	"gade/srv-gade-point/metrics"
	"gade/srv-gade-point/models"
)

type psqlMetricRepository struct {
	Conn *sql.DB
}

// NewPsqlMetricRepository will create an object that represent the metric.Repository interface
func NewPsqlMetricRepository(Conn *sql.DB) metrics.Repository {
	return &psqlMetricRepository{Conn}
}

func (m *psqlMetricRepository) FindMetric(module string) (string, error) {
	var cacing string

	query := `SELECT module FROM metrics where "module" = $1`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		return "", err
	}

	err = stmt.QueryRow(module).Scan(&cacing)

	if err != nil {
		return "", err
	}

	return cacing, err
}

func (m *psqlMetricRepository) CreateMetric(module string) error {
	var lastID int64
	status := int8(1)

	query := `INSERT INTO metrics ("module", counter) VALUES ($1, $2)  RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		return err
	}

	err = stmt.QueryRow(module, status).Scan(&lastID)

	if err != nil {
		return err
	}

	return nil
}

func (m *psqlMetricRepository) UpdateMetric(module string) error {
	var lastID int64

	query := `UPDATE metrics SET counter = counter + 1 WHERE module = $1 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		return err
	}

	err = stmt.QueryRow(module).Scan(&lastID)

	if err != nil {
		return err
	}

	return nil
}

func (m *psqlMetricRepository) BalanceMetric(counter int64, module string) error {
	var lastID int64

	query := `UPDATE metrics SET counter = counter - $1 WHERE module = $2 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		return err
	}

	err = stmt.QueryRow(counter, module).Scan(&lastID)

	if err != nil {
		return err
	}

	return nil
}

func (m *psqlMetricRepository) GetMetric() ([]models.Metrics, error) {
	var result []models.Metrics

	query := `SELECT ID, Module, counter FROM metrics where counter > 0;`

	rows, err := m.Conn.Query(query)

	if err != nil {

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		metric := models.Metrics{}

		err = rows.Scan(
			&metric.ID,
			&metric.Module,
			&metric.Counter,
		)

		if err != nil {

			return nil, err
		}

		result = append(result, metric)

	}

	return result, nil
}
