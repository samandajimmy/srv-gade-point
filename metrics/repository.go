package metrics

import (
	"gade/srv-gade-point/models"
)

// Repository represent the metric's repository contract
type Repository interface {
	FindMetric(string) (string, error)
	CreateMetric(string) error
	UpdateMetric(string) error
	BalanceMetric(int64, string) error
	GetMetric() ([]models.Metrics, error)
}
