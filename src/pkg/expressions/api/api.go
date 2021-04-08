package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type Api interface {
	Score(url string) int
}

type api struct {
}

func NewApi() Api {
	return &api{}
}

func (a *api) Score(url string) int {
	min := 0
	max := 100
	return rand.Intn(max-min+1) + min
}

const healthCheckPath = "/api/v1/health"
const arithmeticExpressionPath = "/api/v1/arithmetic"

type expression struct {
	expr   string
	result float64
	score  int
}

const healthScore = 5

var expressions = []expression{
	{"2 + 3", 5, 2},
	{"2 + 3 - 1", 4, 1},
	{"2+3-1", 4, 2},
	{"2+3-1*2", 3, 5},
	{"2+3-1*2", 3, 5},
	{"2+3-1*2/2", 4, 5},
	{"1 - 2*4+3-1*2/2", 4, 5},
	{"1-2*4+3-1*2/2", 4, 5},
	{"1-2*4+(3-1)*2/2", 4, 5},
}

func CheckEndpoint(url string) int {
	err := healthCheck(url)
	if err != nil {
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
		return false
	}

	if r.StatusCode != http.StatusOK {
		return false
	}

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}

	return string(b) == fmt.Sprintf("%v", res)
}

func healthCheck(host string) error {
	r, err := http.Get(host + healthCheckPath)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return errors.New("status code is not ")
	}

	return nil
}
