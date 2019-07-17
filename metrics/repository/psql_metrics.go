package repository

import (
	"database/sql"
	"gade/srv-gade-point/metrics"
	"time"

	"github.com/sirupsen/logrus"
)

type psqlMetricRepository struct {
	Conn *sql.DB
}

// NewPsqlMetricRepository will create an object that represent the metric.Repository interface
func NewPsqlMetricRepository(Conn *sql.DB) metrics.Repository {
	return &psqlMetricRepository{Conn}
}

func (m *psqlMetricRepository) FindMetric(job string) (string, error) {
	var lastID string

	query := `SELECT id FROM metrics where job = $1`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)

		return "", err
	}

	err = stmt.QueryRow(&job).Scan(&lastID)

	if err != nil {
		logrus.Debug(err)

		return "", err
	}

	return lastID, err
}

func (m *psqlMetricRepository) CreateMetric(job string) error {
	var lastID int64
	counter := int8(1)
	status := int8(0)
	now := time.Now()

	query := `INSERT INTO metrics (job, counter, status, creation_time) VALUES ($1, $2, $3, $4) RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	err = stmt.QueryRow(&job, &counter, &status, &now).Scan(&lastID)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	return nil
}

func (m *psqlMetricRepository) UpdateMetric(job string) error {
	var lastID int64
	now := time.Now()

	query := `UPDATE metrics SET counter = counter + 1, modification_time = $1 WHERE job = $2 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	err = stmt.QueryRow(&now, &job).Scan(&lastID)

	if err != nil {
		logrus.Debug(err)

		return err
	}

	return nil
}
