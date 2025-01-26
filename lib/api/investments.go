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

var ProductionInvestmentsServer = InvestmentsServer{CalculateGeneralPresentValue: i.CalculatePresentValue}

func (pvs *InvestmentsServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	pvs.route(responseWriter, request)
}

func (pvs *InvestmentsServer) route(responseWriter http.ResponseWriter, request *http.Request) {
	subPath := strings.TrimPrefix(request.URL.Path, UrlPrefix)
	if strings.HasPrefix(subPath, "/general-present-value") {
		pvs.handleGeneralPresentValue(responseWriter, request)
	} else {
		pvs.handleUnknownPath(responseWriter, request)
	}
}

func (pvs *InvestmentsServer) handleInappropriateMethod(responseWriter http.ResponseWriter) {
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}

func (pvs *InvestmentsServer) handleUnknownPath(responseWriter http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	responseWriter.WriteHeader(http.StatusNotFound)
	responseWriter.Write([]byte(fmt.Sprintf("unknown path %q", path)))
}

type GeneralPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type GeneralPresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}

func (pvs *InvestmentsServer) handleGeneralPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		pvs.handleInappropriateMethod(responseWriter)
	}
	var body GeneralPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := pvs.CalculateGeneralPresentValue(body.Periods...)
	WriteResponseBody(responseWriter, GeneralPresentValueResponseBody{presentValue})
}
