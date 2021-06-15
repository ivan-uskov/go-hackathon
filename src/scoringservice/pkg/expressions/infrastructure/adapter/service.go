package adapter

import (
	"bytes"
	"errors"
	"go-hackathon/src/common/infrastructure"
	"go-hackathon/src/scoringservice/pkg/expressions/application/adapter"
	"io/ioutil"
	"net/http"
	"strconv"
)

const healthCheckPath = "/api/v1/health"
const arithmeticExpressionPath = "/api/v1/arithmetic"

type externalServiceAdapter struct {
	endpoint string
}

func (e *externalServiceAdapter) HealthCheck() bool {
	r, err := http.Head(e.endpoint + healthCheckPath)
	if err != nil {
		return false
	}

	if r.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func (e *externalServiceAdapter) Calculate(expr string) (float64, error) {
	r, err := http.Post(e.endpoint+arithmeticExpressionPath, "text/plain", bytes.NewBuffer([]byte(expr)))
	if err != nil {
		return 0, err
	}

	if r.StatusCode != http.StatusOK {
		return 0, errors.New("status code is not 200")
	}

	defer infrastructure.Close(r.Body)
	clientRes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(string(clientRes), 64)
}

func NewExternalServiceAdapterFactory() adapter.ExternalServiceAdapterFactory {
	return &externalServiceAdapterFactory{}
}

type externalServiceAdapterFactory struct{}

func (e *externalServiceAdapterFactory) Create(endpoint string) adapter.ExternalServiceAdapter {
	return &externalServiceAdapter{endpoint}
}
