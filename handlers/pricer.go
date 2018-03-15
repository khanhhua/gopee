package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.axa.com/axa-singapore-meetups/gopee/engine"
	"github.axa.com/axa-singapore-meetups/gopee/libhttp"
	"github.com/tealeg/xlsx"
	"github.com/xuri/efp"
)

type Page struct {
	Result string
}

func GetPricerContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/pricer-content.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	excelFileName := "/Users/khanh/Downloads/dup.xlsx"
	fmt.Printf("Reading XLSX...\n")
	xlFile, err := xlsx.OpenFile(excelFileName)

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	ps := efp.ExcelParser()
	Formula := `=IF(OR(CalculatorNB!$B$12="Decline",CalculatorNB!$B$12="Refer"),CalculatorNB!$B$12,CalculatorNB!E48)`
	ps.Parse(Formula)

	fmt.Printf("%s\n", Formula)
	for _, token := range ps.Tokens.Items {
		fmt.Printf("- %s %s %s\n", token.TValue, token.TType, token.TSubType)
	}

	formulaEngine := engine.New(&ps)
	result := formulaEngine.Execute(xlFile)
	page := Page{result}

	fmt.Printf("Result: %s\n", result)

	tmpl.Execute(w, page)
}
