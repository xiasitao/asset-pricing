package api

import (
	"net/http"

	i "xiasitao.de/asset-pricing/lib/investments"
)

type InvestmentsRouter struct {
	CalculateGeneralPresentValue func(...i.Period) i.CurrencyUnit
	http.Handler
}

func NewInvestmentsRouter(calculateGeneralPresentValue func(...i.Period) i.CurrencyUnit) *InvestmentsRouter {
	router := new(InvestmentsRouter)
	router.CalculateGeneralPresentValue = calculateGeneralPresentValue
	mux := http.NewServeMux()
	mux.HandleFunc("/general-present-value", router.handleGeneralPresentValue)
	router.Handler = mux
	return router
}

var ProductionInvestmentsRouter = NewInvestmentsRouter(i.CalculatePresentValue)

type GeneralPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type GeneralPresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}

func (pvs *InvestmentsRouter) handleGeneralPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body GeneralPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := pvs.CalculateGeneralPresentValue(body.Periods...)
	WriteResponseBody(responseWriter, GeneralPresentValueResponseBody{presentValue})
}
