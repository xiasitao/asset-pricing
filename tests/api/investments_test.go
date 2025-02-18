package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	a "xiasitao.de/asset-pricing/lib/api"
	i "xiasitao.de/asset-pricing/lib/investments"
)

func postRequestResponseWriterFactory(path string, body string) (request *http.Request, responseWriter *httptest.ResponseRecorder) {
	request, _ = http.NewRequest(http.MethodPost, path, strings.NewReader(body))
	responseWriter = httptest.NewRecorder()
	return
}

const mockPresentValue i.CurrencyUnit = 123.5

var expectedResponseFromMockHandlers = a.PresentValueResponseBody{PresentValue: i.CurrencyUnit(mockPresentValue)}

var mockHandlers a.InvestmentsRouterHandlers = a.InvestmentsRouterHandlers{
	PerpetuityPresentValue: func(cashflow i.CurrencyUnit, interest float64) i.CurrencyUnit {
		return mockPresentValue
	},
	LumpSumPresentValue: func(lumpSum i.CurrencyUnit, interest float64, periods int) i.CurrencyUnit {
		return mockPresentValue
	},
	AnnuityPresentValue: func(cashflow i.CurrencyUnit, interest float64, periods int) i.CurrencyUnit {
		return mockPresentValue
	},
	GeneralFiniteCashflowPresentValue: func(...i.Period) i.CurrencyUnit {
		return mockPresentValue
	},
}

const perpetuityPresentValueEndpoint = "/perpetuity-present-value"
const lumpSumPresentValueEndpoint = "/lump-sum-present-value"
const annuityPresentValueEndpoint = "/annuity-present-value"
const generalFiniteCashflowPresentValueEndpoint = "/general-finite-cashflow-present-value"

func TestInvestmentsRouter(t *testing.T) {
	t.Run("test unknown path", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/unknown-path", "")
		router := a.NewInvestmentsRouter("", mockHandlers)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusNotFound)
	})
}

func TestPresentValueWithMock(t *testing.T) {
	t.Run("path prefix", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/prefix"+generalFiniteCashflowPresentValueEndpoint, "{}")
		router := a.NewInvestmentsRouter("/prefix", mockHandlers)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusOK)
	})

	t.Run("malformed request body", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory(generalFiniteCashflowPresentValueEndpoint, "{malformed}")
		router := a.NewInvestmentsRouter("", mockHandlers)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusUnprocessableEntity)
	})

	testMethod := func(t *testing.T, endpoint string) {
		t.Helper()
		request, _ := http.NewRequest(http.MethodGet, endpoint, strings.NewReader(""))
		responseWriter := httptest.NewRecorder()
		server := a.NewInvestmentsRouter("", mockHandlers)
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusMethodNotAllowed)
	}

	testEndpointResponse := func(t *testing.T, endpoint string, requestBody any) {
		t.Helper()
		requestBodyBuffer := strings.Builder{}
		json.NewEncoder(&requestBodyBuffer).Encode(requestBody)
		request, responseWriter := postRequestResponseWriterFactory(
			endpoint,
			requestBodyBuffer.String(),
		)
		router := a.NewInvestmentsRouter("", mockHandlers)
		router.ServeHTTP(responseWriter, request)
		got := a.PresentValueResponseBody{}
		json.NewDecoder(responseWriter.Body).Decode(&got)
		expected := expectedResponseFromMockHandlers
		assertResponseBody(t, got, expected)
	}

	t.Run("perpetuity", func(t *testing.T) {
		testMethod(t, perpetuityPresentValueEndpoint)
		testEndpointResponse(t, perpetuityPresentValueEndpoint, a.PerpetuityPresentValueRequestBody{Cashflow: 1.0, Interest: 1.0})
	})

	t.Run("lump sum", func(t *testing.T) {
		testMethod(t, lumpSumPresentValueEndpoint)
		testEndpointResponse(t, lumpSumPresentValueEndpoint, a.LumpSumPresentValueRequestBody{Lump: 1.0, Interest: 1.0, Periods: 5})
	})

	t.Run("annuity", func(t *testing.T) {
		testMethod(t, annuityPresentValueEndpoint)
		testEndpointResponse(t, annuityPresentValueEndpoint, a.AnnuityPresentValueRequestBody{Cashflow: 1.0, Interest: 1.0, Periods: 5})
	})

	t.Run("general finite cashflow", func(t *testing.T) {
		testMethod(t, generalFiniteCashflowPresentValueEndpoint)
		testEndpointResponse(t, generalFiniteCashflowPresentValueEndpoint,
			a.GeneralFiniteCashflowPresentValueRequestBody{Periods: []i.Period{{Cashflow: 1.0, Interest: 1.0}}},
		)
	})

}

func TestProductionInvestmentsRouterIntegration(t *testing.T) {
	testEndpointResponse := func(t *testing.T, endpoint string, requestBody any, expectedResponseBody a.PresentValueResponseBody) {
		t.Helper()
		requestBodyBuffer := strings.Builder{}
		json.NewEncoder(&requestBodyBuffer).Encode(requestBody)

		request, responseWriter := postRequestResponseWriterFactory(
			"/investments"+endpoint,
			requestBodyBuffer.String(),
		)
		a.ProductionInvestmentsRouter.ServeHTTP(responseWriter, request)

		got := a.PresentValueResponseBody{}
		json.NewDecoder(responseWriter.Body).Decode(&got)
		assertStatus(t, responseWriter, http.StatusOK)
		assertResponseBody(t, got, expectedResponseBody)
	}

	t.Run("perpetuity", func(t *testing.T) {
		requestBody := a.PerpetuityPresentValueRequestBody{
			Cashflow: 1.0, Interest: 0.2,
		}
		expectedResponseBody := a.PresentValueResponseBody{
			PresentValue: 5.0,
		}
		testEndpointResponse(t, perpetuityPresentValueEndpoint, requestBody, expectedResponseBody)
	})

	t.Run("lump sum", func(t *testing.T) {
		requestBody := a.LumpSumPresentValueRequestBody{
			Lump:     1000.0,
			Interest: 0.2,
			Periods:  5,
		}
		expectedResponseBody := a.PresentValueResponseBody{
			PresentValue: 1000.0 / 1.2 / 1.2 / 1.2 / 1.2 / 1.2,
		}
		testEndpointResponse(t, lumpSumPresentValueEndpoint, requestBody, expectedResponseBody)
	})

	t.Run("annuity", func(t *testing.T) {
		requestBody := a.AnnuityPresentValueRequestBody{
			Cashflow: 1.0,
			Interest: 0.2,
			Periods:  5,
		}
		expectedResponseBody := a.PresentValueResponseBody{
			PresentValue: 1.0 * (1 - 1/1.2/1.2/1.2/1.2/1.2) / 0.2,
		}
		testEndpointResponse(t, annuityPresentValueEndpoint, requestBody, expectedResponseBody)
	})

	t.Run("general finite cashflow", func(t *testing.T) {
		requestBody := a.GeneralFiniteCashflowPresentValueRequestBody{Periods: []i.Period{{Cashflow: 1.0, Interest: 1.0}}}
		expectedResponseBody := a.PresentValueResponseBody{PresentValue: 0.5}
		testEndpointResponse(t, generalFiniteCashflowPresentValueEndpoint, requestBody, expectedResponseBody)
	})

}

func assertStatus(t testing.TB, responseWriter *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := responseWriter.Code

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func assertResponseBody(t testing.TB, got, expected a.PresentValueResponseBody) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}
