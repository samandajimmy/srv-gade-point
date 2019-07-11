package metrics

// Repository represent the metric's repository contract
type Repository interface {
	FindMetric(string) (string, error)
	CreateMetric(string) error
	UpdateMetric(string) error
}
