package api

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
)

const healthCheckPath = "/api/v1/health"
const arithmeticExpressionPath = "/api/v1/arithmetic"

type expression struct {
	expr   string
	result float64
	score  int
}

const healthScore = 5

var expressions = []expression{
	//check "+", "-", spaces and more than two args
	{"1", 1, 1},
	{"2+2", 4, 1},
	{"2+2-4", 0, 1},
	{"2 + 3", 5, 1},
	{"8+ 5", 13, 1},
	{"12 + 5", 17, 1},
	{"2 + 3 - 1", 4, 1},
	{"2+ 3+4+ 5 + 6", 20, 1},
	//check "*", "/" tens, hundreds...
	{"2*3", 6, 2},
	{"8/2", 4, 2},
	{"10*3", 30, 2},
	{"448/8", 56, 2},
	{"2*3*1*2*2", 24, 2},
	{"500/25/20", 1, 2},
	{"500*250/5", 25000, 2},
	{"4* 6/ 2 *10 / 5", 24, 2},
	//check negative and fractional numbers
	{"-2", -2, 3},
	{"-8+3", -5, 3},
	{"-13/-3", 4.33, 3},
	//check priorities and brackets
	{"1-3*3", -8, 4},
	{"1 - 2*4+3-1*2/2", -5, 4},
	{"1-2*4+3-(1*2/2)", -5, 4},
	{"1-2*(4+2-1)*(2/2)", -9, 4},
	{"(((((8/2)-3)*4)+8)*9)+0", 108, 4},
	{"(((((288/2)))))", 144, 4},
	{"(1*(0))+19", 19, 4},
	{"-2 -(9 *6)", -56, 4},
	{"(-2 )* 27", -54, 4},
	{"(4+6-(-2) *6)/92", 0.24, 4},
	//check boundary values
	{fmt.Sprintf("%v+%v-%v", math.MaxInt64, math.MaxInt64, math.MaxInt64), math.MaxInt64, 5},
	{fmt.Sprintf("%v+%v-%v", math.MaxFloat64, math.MaxFloat64, math.MaxFloat64), math.MaxFloat64, 5},
}

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
