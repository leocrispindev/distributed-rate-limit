package domain

type MetricsRepository interface {
	CountMetric(clientName string) error
}
