package usecase

import (
	"gade/srv-gade-point/metrics"
	"gade/srv-gade-point/models"
	"time"

	"github.com/sirupsen/logrus"
)

type metricUseCase struct {
	metricRepo     metrics.Repository
	contextTimeout time.Duration
}

// NewMetricUseCase  will create new an metricUseCase object representation of metrics.UseCase interface
func NewMetricUseCase(met metrics.Repository, timeout time.Duration) metrics.UseCase {
	return &metricUseCase{
		metricRepo:     met,
		contextTimeout: timeout,
	}
}

func (met *metricUseCase) AddMetric(module string) error {

	data, err := met.metricRepo.FindMetric(module)
	if data == "" {
		err = met.metricRepo.CreateMetric(module)

		if err != nil {
			return models.ErrCreateMetric
		}
	} else {
		err = met.metricRepo.UpdateMetric(module)

		if err != nil {
			return models.ErrUpdateMetric
		}
	}

	return err
}

func (met *metricUseCase) SendMetric() error {

	list, err := met.metricRepo.GetMetric()

	if list != nil {

		logrus.Debug("Start Send Metric to Prometheus")
		for i := 0; i < len(list); i++ {

			var counter = list[i].Counter
			var module = list[i].Module

			err = met.metricRepo.BalanceMetric(*counter, module)

			if err != nil {
				return models.ErrUpdateMetric
			}

		}
		logrus.Debug("End Send Metric to Prometheus")
	}

	return err
}
