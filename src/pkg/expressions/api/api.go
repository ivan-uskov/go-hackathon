package api

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/transport"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

const healthCheckPath = "/api/v1/health"
const arithmeticExpressionPath = "/api/v1/arithmetic"

const epsilon = 0.01
const healthScore = 5

type Api interface {
	Score(url string) int
}

type api struct{}

func NewApi() Api {
	return &api{}
}

func (a *api) Score(url string) int {
	err := healthCheck(url)
	if err != nil {
		log.Debug(err)
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
		log.Debug(err)
		return false
	}

	if r.StatusCode != http.StatusOK {
		return false
	}

	defer transport.CloseBody(r.Body)
	clientRes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debug(err)
		return false
	}

	clientNum, err := strconv.ParseFloat(string(clientRes), 64)
	if err != nil {
		log.Debug(err)
		return false
	}

	return compareFloats(clientNum, res)
}

func healthCheck(host string) error {
	r, err := http.Head(host + healthCheckPath)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New("status code is not 200")
	}

	return nil
}

func compareFloats(left float64, right float64) bool {
	return math.Abs(left-right) <= epsilon
}
