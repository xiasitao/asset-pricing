package api

import (
	"net/http"

	i "xiasitao.de/asset-pricing/lib/investments"
)

type InvestmentsRouter struct {
	CalculateGeneralPresentValue func(...i.Period) i.CurrencyUnit
	http.Handler
}

func NewInvestmentsRouter(prefix string, calculateGeneralFiniteCashflowPresentValue func(...i.Period) i.CurrencyUnit) *InvestmentsRouter {
	router := InvestmentsRouter{}
	router.CalculateGeneralPresentValue = calculateGeneralFiniteCashflowPresentValue
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/general-finite-cashflow-present-value", router.handleGeneralFiniteCashflowPresentValue)
	router.Handler = mux
	return &router
}

var ProductionInvestmentsRouter = NewInvestmentsRouter("/investments", i.CalculateGeneralFiniteCashflowPresentValue)

type GeneralFiniteCashflowPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type GeneralFiniteCashflowPresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}

func (pvs *InvestmentsRouter) handleGeneralFiniteCashflowPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body GeneralFiniteCashflowPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := pvs.CalculateGeneralPresentValue(body.Periods...)
	WriteResponseBody(responseWriter, GeneralFiniteCashflowPresentValueResponseBody{presentValue})
}
