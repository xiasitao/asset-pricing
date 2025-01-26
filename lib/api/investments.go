package api

import (
	"encoding/json"
	"net/http"
	"strings"

	i "xiasitao.de/asset-pricing/lib/investments"
)

type PresentValueServer struct {
	CalculateGeneralPresentValue func(...i.Period) i.CurrencyUnit
}

func (pvs *PresentValueServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	pvs.route(responseWriter, request)
}

func (pvs *PresentValueServer) route(responseWriter http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if strings.HasPrefix(path, "/general-present-value") {
		pvs.handleGeneralPresentValue(responseWriter, request)
	} else {

	}
}

func (pvs *PresentValueServer) handleGeneralPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	var body GeneralPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}

	presentValue := pvs.CalculateGeneralPresentValue(body.Periods...)
	json.NewEncoder(responseWriter).Encode(GeneralPresentValueResponseBody{presentValue})
}

type GeneralPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type GeneralPresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}
