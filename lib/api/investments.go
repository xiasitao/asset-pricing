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
	LumpSumPresentValue               func(lump i.CurrencyUnit, interest float64, periods int) i.CurrencyUnit
	AnnuityPresentValue               func(cashflow i.CurrencyUnit, interest float64, periods int) i.CurrencyUnit
	GeneralFiniteCashflowPresentValue func(...i.Period) i.CurrencyUnit
}

func NewInvestmentsRouter(prefix string, handlers InvestmentsRouterHandlers) *InvestmentsRouter {
	router := InvestmentsRouter{}
	router.handlers = handlers
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/perpetuity-present-value", router.handlePerpetuityPresentValue)
	mux.HandleFunc(prefix+"/lump-sum-present-value", router.handleLumpSumPresentValue)
	mux.HandleFunc(prefix+"/annuity-present-value", router.handleAnnuityPresentValue)
	mux.HandleFunc(prefix+"/general-finite-cashflow-present-value", router.handleGeneralFiniteCashflowPresentValue)
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

type LumpSumPresentValueRequestBody struct {
	Lump     i.CurrencyUnit `json:"lump"`
	Interest float64        `json:"interest"`
	Periods  int            `json:"periods"`
}

func (ir *InvestmentsRouter) handleLumpSumPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body LumpSumPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := ir.handlers.LumpSumPresentValue(body.Lump, body.Interest, body.Periods)
	WriteResponseBody(responseWriter, PresentValueResponseBody{presentValue})
}

type AnnuityPresentValueRequestBody struct {
	Cashflow i.CurrencyUnit
	Interest float64
	Periods  int
}

func (ir *InvestmentsRouter) handleAnnuityPresentValue(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		responseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body AnnuityPresentValueRequestBody
	err := ReadRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}
	presentValue := ir.handlers.AnnuityPresentValue(body.Cashflow, body.Interest, body.Periods)
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
