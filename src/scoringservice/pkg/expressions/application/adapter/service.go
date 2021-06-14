package adapter

type ExternalServiceAdapterFactory interface {
	Create(endpoint string) ExternalServiceAdapter
}

type ExternalServiceAdapter interface {
	HealthCheck() bool
	Calculate(expr string) (float64, error)
}
