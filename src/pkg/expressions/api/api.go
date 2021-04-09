package api

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const healthCheckPath = "/api/v1/health"
const arithmeticExpressionPath = "/api/v1/arithmetic"

const healthScore = 5

type Api interface {
	Score(url string) int
}

type api struct {
}

func NewApi() Api {
	return &api{}
}

func (a *api) Score(url string) int {
	err := healthCheck(url)
	if err != nil {
		log.Error(err)
		return 0
	}

	score := healthScore

	for _, e := range expressions {
		if arithmetic(url, e.expr, e.result) {
			score += e.score
		}
	}

	return score
}

func arithmetic(host string, expr string, res float64) bool {
	r, err := http.Post(host+arithmeticExpressionPath, "text/plain", bytes.NewBuffer([]byte(expr)))
	if err != nil {
		log.Error(err)
		return false
	}

	if r.StatusCode != http.StatusOK {
		return false
	}

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return false
	}

	return string(b) == fmt.Sprintf("%v", res)
}

func healthCheck(host string) error {
	r, err := http.Head(host + healthCheckPath)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New("status code is not ")
	}

	return nil
}
