package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"io/ioutil"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	dbxFiles "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/gorilla/mux"
	"github.com/khanhhua/gopee/dao"

	f1Engine "github.com/khanhhua/formula1/engine"
	xlsx "github.com/tealeg/xlsx"
)

type param struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type fun struct {
	ID             int64   `json:"id"`
	FnName         string  `json:"fnName"`
	XlsxName       string  `json:"xlsxName"`
	InputMappings  []param `json:"inputMappings"`
	OutputMappings []param `json:"outputMappings"`
}

var payload struct {
	Fun fun `json:"fun"`
}

// Compose composes function
func Compose(w http.ResponseWriter, r *http.Request) {
	var err error
	var data []byte
	clientKey := r.Header.Get("x-client-key")
	if len(clientKey) == 0 {
		http.Error(w, "Not authorized", 403)
		return
	}

	if data, err = ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if err = json.Unmarshal(data, &payload); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var id int64
	f := payload.Fun
	if id, err = dao.CreateFunc(clientKey, f.FnName, f.XlsxName,
		paramDeclarationsToMappings(f.InputMappings),
		paramDeclarationsToMappings(f.OutputMappings)); err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		payload.Fun.ID = id
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func paramDeclarationsToMappings(params []param) map[string]string {
	ret := make(map[string]string)
	for _, p := range params {
		ret[p.Name] = p.Address
	}

	return ret
}

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

	var rawInputs map[string]interface{}
	if fun, err := dao.GetFuncByName(clientKey, fnName); err != nil {
		fmt.Printf("Function '%s' not found:", fnName)
		http.Error(w, "Bad request", 404)
		return
	} else if data, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, "Bad request", 400)
		return
	} else {
		json.Unmarshal(data, &rawInputs)
		inputs := make(map[string]string)
		for key, input := range rawInputs {
			switch input.(type) {
			case string:
				inputs[key] = input.(string)
				break
			case bool:
				inputs[key] = strconv.FormatBool(input.(bool))
				break
			case float64:
				inputs[key] = strconv.FormatFloat(input.(float64), 'E', -1, 64)
				break
			}
		}

		if outputs, err := execute(fun, inputs); err != nil {
			fmt.Printf("Error during executing: %v", err)
			http.Error(w, "Execution error", 500)

			return
		} else {
			encoder := json.NewEncoder(w)
			encoder.Encode(outputs)
		}
	}
}

func execute(fun dao.FuncSpec, inputs map[string]string) (outputs map[string]string, err error) {
	config := dropbox.Config{
		Token:    fun.DropboxAccessToken,
		LogLevel: dropbox.LogInfo, // if needed, set the desired logging level. Default is off
	}
	dbxClient := dbxFiles.New(config)
	filesDownloadArgs := dbxFiles.NewDownloadArg(fun.XlsxFile)
	if _, fileContent, dbxErr := dbxClient.Download(filesDownloadArgs); dbxErr != nil {
		fmt.Printf("Could not download.\n %v", dbxErr)
		err = dbxErr
		return
	} else if bs, fileErr := ioutil.ReadAll(fileContent); fileErr != nil {
		fmt.Printf("Could not read.\n %v", fileErr)

		err = fileErr
		return
	} else if xlFile, xlsxErr := xlsx.OpenBinary(bs); xlsxErr != nil {
		err = xlsxErr
		return
	} else {
		fmt.Printf("File content: %v", bs)

		mappedInputs := make(map[string]string)
		for param, address := range fun.InputMappings {
			mappedInputs[address] = inputs[param]
		}

		mappedOutputs := make(map[string]string)
		for _, address := range fun.OutputMappings {
			mappedOutputs[address] = ""
		}

		engine := f1Engine.NewEngine(xlFile)
		if err = engine.Execute(mappedInputs, &mappedOutputs); err == nil {
			outputs = make(map[string]string)
			for param, address := range fun.OutputMappings {
				outputs[param] = mappedOutputs[address]
			}
		}
	}

	return
}
