package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khanhhua/gopee/dao"
)
import "io/ioutil"

// Call calls a precomposed function
func Call(w http.ResponseWriter, r *http.Request) {
	clientKey := r.Header.Get("x-client-key")
	if len(clientKey) == 0 {
		http.Error(w, "Not authorized", 403)
		return
	}

	fnName := mux.Vars(r)["fnName"]
	if len(fnName) == 0 {
		http.Error(w, "Bad request", 400)
		return
	}

	var inputs map[string]interface{}
	if fun, err := dao.GetFuncByName(clientKey, fnName); err != nil {
		http.Error(w, "Bad request", 404)
		return
	} else if data, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, "Bad request", 404)
		return
	} else {
		json.Unmarshal(data, &inputs)
		if output, err := execute(fun, inputs); err != nil {
			http.Error(w, "Execution error", 500)
			return
		} else {
			encoder := json.NewEncoder(w)
			encoder.Encode(output)
		}
	}
}

func execute(fun dao.FuncSpec, inputs map[string]interface{}) (output map[string]interface{}, err error) {
	output = make(map[string]interface{})
	output["A1"] = "OK"
	output["A2"] = inputs["A2"]

	return
}
