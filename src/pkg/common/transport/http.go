package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go-hackaton/src/pkg/common/application/errors"
	"io/ioutil"
	"net/http"
)

func ProcessError(w http.ResponseWriter, e error) {
	if e == errors.InternalError {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	} else {
		http.Error(w, e.Error(), http.StatusBadRequest)
	}
}

func RenderJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error(err)
		ProcessError(w, errors.InternalError)
		return
	}
}

func ReadJson(r *http.Request, output interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer func() {
		log.Error(r.Body.Close())
	}()

	err = json.Unmarshal(b, &output)
	if err != nil {
		err = fmt.Errorf("can't parse %s to json", b)
	}

	return err
}

func Parameter(r *http.Request, key string) (string, bool) {
	val, found := mux.Vars(r)["ID"]
	return val, found
}