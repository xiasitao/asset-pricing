package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"xiasitao.de/asset-pricing/lib/api"
	"xiasitao.de/asset-pricing/lib/investments"
)

func postRequestResponseWriterFactory(path string, body string) (request *http.Request, responseWriter *httptest.ResponseRecorder) {
	request, _ = http.NewRequest(http.MethodPost, path, strings.NewReader(body))
	responseWriter = httptest.NewRecorder()
	return
}

const mockPresentValue investments.CurrencyUnit = 123.5

var expectedResponseFromMockCalculateGeneralFiniteCashflowPresentValue = api.PresentValueResponseBody{PresentValue: investments.CurrencyUnit(mockPresentValue)}

func mockCalculateGeneralFiniteCashflowPresentValue(...investments.Period) investments.CurrencyUnit {
	return mockPresentValue
}

const generalPresentValueEndpoint = "/general-finite-cashflow-present-value"

func TestInvestmentsRouter(t *testing.T) {
	t.Run("test unknown path", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/unknown-path", "")
		router := api.NewInvestmentsRouter("", mockCalculateGeneralFiniteCashflowPresentValue)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusNotFound)
	})
}

func TestPresentValueWithMock(t *testing.T) {
	t.Run("path prefix", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/prefix"+generalPresentValueEndpoint, "{}")
		router := api.NewInvestmentsRouter("/prefix", mockCalculateGeneralFiniteCashflowPresentValue)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusOK)
	})

	t.Run("malformed request body", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory(generalPresentValueEndpoint, "{malformed}")
		router := api.NewInvestmentsRouter("", mockCalculateGeneralFiniteCashflowPresentValue)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusUnprocessableEntity)
	})

	t.Run("wrong method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, generalPresentValueEndpoint, strings.NewReader(""))
		responseWriter := httptest.NewRecorder()
		server := api.NewInvestmentsRouter("", mockCalculateGeneralFiniteCashflowPresentValue)
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusMethodNotAllowed)
	})

	t.Run("general finite cashflow", func(t *testing.T) {
		requestBodyBuffer := strings.Builder{}
		json.NewEncoder(&requestBodyBuffer).Encode(api.GeneralFiniteCashflowPresentValueRequestBody{Periods: []investments.Period{{Cashflow: 1.0, Interest: 1.0}}})
		request, responseWriter := postRequestResponseWriterFactory(
			generalPresentValueEndpoint,
			requestBodyBuffer.String(),
		)
		router := api.NewInvestmentsRouter("", mockCalculateGeneralFiniteCashflowPresentValue)
		router.ServeHTTP(responseWriter, request)

		got := api.PresentValueResponseBody{}
		json.NewDecoder(responseWriter.Body).Decode(&got)
		expected := expectedResponseFromMockCalculateGeneralFiniteCashflowPresentValue
		assertResponseBody(t, got, expected)
	})

}

func TestPresentValueIntegration(t *testing.T) {
	requestBodyBuffer := strings.Builder{}
	json.NewEncoder(&requestBodyBuffer).Encode(api.GeneralFiniteCashflowPresentValueRequestBody{Periods: []investments.Period{{Cashflow: 1.0, Interest: 1.0}}})

	request, responseWriter := postRequestResponseWriterFactory(
		"/investments"+generalPresentValueEndpoint,
		requestBodyBuffer.String(),
	)
	api.ProductionInvestmentsRouter.ServeHTTP(responseWriter, request)

	got := api.PresentValueResponseBody{}
	json.NewDecoder(responseWriter.Body).Decode(&got)
	expected := api.PresentValueResponseBody{PresentValue: 0.5}
	assertStatus(t, responseWriter, http.StatusOK)
	assertResponseBody(t, got, expected)
}
func assertStatus(t testing.TB, responseWriter *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := responseWriter.Code

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func assertResponseBody(t testing.TB, got, expected api.PresentValueResponseBody) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}
