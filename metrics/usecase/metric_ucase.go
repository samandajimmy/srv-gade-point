package usecase

import (
	"gade/srv-gade-point/metrics"
	"gade/srv-gade-point/models"
	"time"
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

func (met *metricUseCase) AddMetric(job string) error {

	data, err := met.metricRepo.FindMetric(job)

	if data == "" {
		err = met.metricRepo.CreateMetric(job)

		if err != nil {
			return models.ErrCreateMetric
		}
	} else {
		err = met.metricRepo.UpdateMetric(job)

		if err != nil {
			return models.ErrUpdateMetric
		}
	}

	return err
}
