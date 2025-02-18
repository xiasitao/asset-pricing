package api

import (
	"net/http"

	i "xiasitao.de/asset-pricing/lib/investments"
)

type InvestmentsRouter struct {
	handlers InvestmentsRouterHandlers
	http.Handler
}

type InvestmentsRouterHandlers struct {
	PerpetuityPresentValue            func(cashflow i.CurrencyUnit, interest float64) i.CurrencyUnit
	GeneralFiniteCashflowPresentValue func(...i.Period) i.CurrencyUnit
}

func NewInvestmentsRouter(prefix string, handlers InvestmentsRouterHandlers) *InvestmentsRouter {
	router := InvestmentsRouter{}
	router.handlers = handlers
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/general-finite-cashflow-present-value", router.handleGeneralFiniteCashflowPresentValue)
	mux.HandleFunc(prefix+"/perpetuity-present-value", router.handlePerpetuityPresentValue)
	router.Handler = mux
	return &router
}

var productionHandlers InvestmentsRouterHandlers = InvestmentsRouterHandlers{
	GeneralFiniteCashflowPresentValue: i.CalculateGeneralFiniteCashflowPresentValue,
}
var ProductionInvestmentsRouter = NewInvestmentsRouter("/investments", productionHandlers)

type PresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}
type PerpetuityPresentValueRequestBody struct {
	Cashflow i.CurrencyUnit `json:"cashflow"`
	Interest float64        `json:"interest"`
}

func (ir *InvestmentsRouter) handlePerpetuityPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body PerpetuityPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := ir.handlers.PerpetuityPresentValue(body.Cashflow, body.Interest)
	WriteResponseBody(responseWriter, PresentValueResponseBody{presentValue})
}

type GeneralFiniteCashflowPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

func (ir *InvestmentsRouter) handleGeneralFiniteCashflowPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body GeneralFiniteCashflowPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := ir.handlers.GeneralFiniteCashflowPresentValue(body.Periods...)
	WriteResponseBody(responseWriter, PresentValueResponseBody{presentValue})
}
