package api

import (
	"fmt"
	"net/http"
	"strings"

	i "xiasitao.de/asset-pricing/lib/investments"
)

type InvestmentsServer struct {
	CalculateGeneralPresentValue func(...i.Period) i.CurrencyUnit
}

func (pvs *InvestmentsServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	pvs.route(responseWriter, request)
}

func (pvs *InvestmentsServer) route(responseWriter http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if strings.HasPrefix(path, "/general-present-value") {
		pvs.handleGeneralPresentValue(responseWriter, request)
	} else {
		pvs.handleUnknownPath(responseWriter, request)
	}
}

func (pvs *InvestmentsServer) handleUnknownPath(responseWriter http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	responseWriter.WriteHeader(404)
	responseWriter.Write([]byte(fmt.Sprintf("unknown path %q", path)))
}

type GeneralPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type GeneralPresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}

func (pvs *InvestmentsServer) handleGeneralPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	var body GeneralPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := pvs.CalculateGeneralPresentValue(body.Periods...)
	WriteResponseBody(responseWriter, GeneralPresentValueResponseBody{presentValue})
}
