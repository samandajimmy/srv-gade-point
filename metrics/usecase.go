package metrics

// UseCase represent the metric's usecases
type UseCase interface {
	AddMetric(string) error
}
