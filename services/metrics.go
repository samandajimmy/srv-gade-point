package services

import (
	"gade/srv-gade-point/metrics"
	"gade/srv-gade-point/models"
)

// MetricService is to handle all metrics service
type MetricService struct {
	metricUsecase metrics.UseCase
}

var some MetricService

// NewMetricHandler is to return a metric service struct
func NewMetricHandler(mu metrics.UseCase) {
	some = MetricService{
		metricUsecase: mu,
	}
}

// AddMetric is to add metric to db
func AddMetric(module string) error {

	err := some.metricUsecase.AddMetric(module)

	if err != nil {
		return models.ErrCreateMetric
	}

	return err
}

// SendMetric is to add metric to promotheus
func SendMetric() error {
	err := some.metricUsecase.SendMetric()

	if err != nil {
		return models.ErrSendMetric
	}

	return err

}
