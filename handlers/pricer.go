package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	g "github.com/khanhhua/formula1/engine"
	"github.com/tealeg/xlsx"
)

type Page struct {
	Result string
}

func ViewPricer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tmpl, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/pricer-content.html.tmpl")
	tmpl.Execute(w, nil)
}

func UploadPricer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("r: %v\n", r)
	if body := r.Body; body == nil {
		w.Write([]byte("error"))
		return
	} else if len := r.ContentLength; len == 0 {
		defer body.Close()
		return
	} else {
		defer body.Close()
		println("Saving to ./data/pricer.xlsx...")
		if pricerFile, err := os.OpenFile("./data/pricer.xlsx", os.O_WRONLY|os.O_CREATE, 0666); err != nil {
			w.Write([]byte("perm_error"))
			return
		} else {
			io.Copy(pricerFile, body)
			defer pricerFile.Close()
		}
	}

	w.Write([]byte("ok"))
}

func GetPricerContent() func(w http.ResponseWriter, r *http.Request) {
	// xlDefaultFile, _ := xlsx.OpenFile("./data/dup.xlsx")
	var xlPricerFile *xlsx.File

	// "Input!E18": "2000000.0",
	// "Input!E20": "0.0",

	// "Input!E35": "",

	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json")
		inputs := map[string]string{}
		outputs := &map[string]string{}

		if body := r.Body; body == nil {
			w.Write([]byte("error"))
			return
		} else if len := r.ContentLength; len == 0 {
			defer body.Close()
			return
		} else {
			defer body.Close()
			buffer := make([]byte, r.ContentLength)
			body.Read(buffer)
			fmt.Printf("%s\n", string(buffer))
			xlPricerFile, err = xlsx.OpenFile("./data/pricer.xlsx")
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			var mapping interface{}

			if err = json.Unmarshal(buffer, &mapping); err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			fmt.Printf("BODY: %v\n\n", mapping)
			aInputs := mapping.(map[string]interface{})["inputs"].([]interface{})

			for _, item := range aInputs {
				sItem := item.(map[string]interface{})
				inputs[sItem["address"].(string)] = sItem["value"].(string)
			}

			aOutputs := mapping.(map[string]interface{})["outputs"].([]interface{})
			for _, item := range aOutputs {
				sItem := item.(map[string]interface{})
				(*outputs)[sItem["address"].(string)] = ""
			}

			var engine *g.Engine

			engine = g.NewEngine(xlPricerFile)
			engine.Execute(inputs, outputs)

			marshalled, _ := json.Marshal(*outputs)
			w.Write([]byte(marshalled))
		}
	}
}
