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
	GeneralFiniteCashflowPresentValue func(...i.Period) i.CurrencyUnit
}

func NewInvestmentsRouter(prefix string, handlers InvestmentsRouterHandlers) *InvestmentsRouter {
	router := InvestmentsRouter{}
	router.handlers = handlers
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/general-finite-cashflow-present-value", router.handleGeneralFiniteCashflowPresentValue)
	router.Handler = mux
	return &router
}

var productionHandlers InvestmentsRouterHandlers = InvestmentsRouterHandlers{
	GeneralFiniteCashflowPresentValue: i.CalculateGeneralFiniteCashflowPresentValue,
}
var ProductionInvestmentsRouter = NewInvestmentsRouter("/investments", productionHandlers)

type GeneralFiniteCashflowPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type PresentValueResponseBody struct {
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
	presentValue := pvs.handlers.GeneralFiniteCashflowPresentValue(body.Periods...)
	WriteResponseBody(responseWriter, PresentValueResponseBody{presentValue})
}
