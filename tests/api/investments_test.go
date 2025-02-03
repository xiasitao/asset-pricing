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

var expectedResponseFromMockCalculateGeneratePresentValue = api.GeneralPresentValueResponseBody{PresentValue: investments.CurrencyUnit(mockPresentValue)}

func mockCalculateGeneralPresentValue(...investments.Period) investments.CurrencyUnit {
	return mockPresentValue
}

const generalPresentValueEndpoint = "/general-present-value"

func TestInvestmentsRouter(t *testing.T) {
	t.Run("test unknown path", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory(api.UrlPrefix+"/unknown-path", "")
		router := api.NewInvestmentsRouter(mockCalculateGeneralPresentValue)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusNotFound)
	})
}

func TestPresentValueWithMock(t *testing.T) {
	t.Run("test correct request", func(t *testing.T) {
		requestBodyBuffer := strings.Builder{}
		json.NewEncoder(&requestBodyBuffer).Encode(api.GeneralPresentValueRequestBody{Periods: []investments.Period{{Cashflow: 1.0, Interest: 1.0}}})
		request, responseWriter := postRequestResponseWriterFactory(
			api.UrlPrefix+generalPresentValueEndpoint,
			requestBodyBuffer.String(),
		)
		router := api.NewInvestmentsRouter(mockCalculateGeneralPresentValue)
		router.ServeHTTP(responseWriter, request)

		got := api.GeneralPresentValueResponseBody{}
		json.NewDecoder(responseWriter.Body).Decode(&got)
		expected := expectedResponseFromMockCalculateGeneratePresentValue
		assertResponseBody(t, got, expected)
	})

	t.Run("test malformed request body", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory(api.UrlPrefix+generalPresentValueEndpoint, "{malformed}")
		router := api.NewInvestmentsRouter(mockCalculateGeneralPresentValue)
		router.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusUnprocessableEntity)
	})

	t.Run("test wrong method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, api.UrlPrefix+generalPresentValueEndpoint, strings.NewReader(""))
		responseWriter := httptest.NewRecorder()
		server := api.NewInvestmentsRouter(mockCalculateGeneralPresentValue)
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusMethodNotAllowed)
	})

}

func TestPresentValueIntegration(t *testing.T) {
	requestBodyBuffer := strings.Builder{}
	json.NewEncoder(&requestBodyBuffer).Encode(api.GeneralPresentValueRequestBody{Periods: []investments.Period{{Cashflow: 1.0, Interest: 1.0}}})

	request, responseWriter := postRequestResponseWriterFactory(
		api.UrlPrefix+generalPresentValueEndpoint,
		requestBodyBuffer.String(),
	)
	api.ProductionInvestmentsRouter.ServeHTTP(responseWriter, request)

	got := api.GeneralPresentValueResponseBody{}
	json.NewDecoder(responseWriter.Body).Decode(&got)
	expected := api.GeneralPresentValueResponseBody{PresentValue: 0.5}
	assertResponseBody(t, got, expected)
}
func assertStatus(t testing.TB, responseWriter *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := responseWriter.Code

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func assertResponseBody(t testing.TB, got, expected api.GeneralPresentValueResponseBody) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}
