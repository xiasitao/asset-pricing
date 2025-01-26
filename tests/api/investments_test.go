package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

var expectedResponseFromMockCalculateGeneratePresentValue = fmt.Sprintf("{\"presentValue\":%.1f}", mockPresentValue)

func mockCalculateGeneralPresentValue(...i.Period) i.CurrencyUnit {
	return mockPresentValue
}

func TestInvestmentsServer(t *testing.T) {
	t.Run("test unknown path", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/unknown-path", "")
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusNotFound)
	})
}

func TestPresentValueWithMock(t *testing.T) {
	t.Run("test correct request", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory(
			"/general-present-value",
			"{\"periods\": [{\"cashflow\": 1.0, \"interest\": 1.0}]}",
		)
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)

		got := responseWriter.Body.String()
		expected := expectedResponseFromMockCalculateGeneratePresentValue
		assertResponseBody(t, got, expected)

	})

	t.Run("test malformed request body", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/general-present-value", "{malformed}")
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusUnprocessableEntity)
	})

	t.Run("test wrong method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/general-present-value", strings.NewReader(""))
		responseWriter := httptest.NewRecorder()
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusMethodNotAllowed)
	})

}

func TestPresentValueIntegration(t *testing.T) {
	server := a.ProductionInvestmentsServer
	request, responseWriter := postRequestResponseWriterFactory(
		"/general-present-value",
		"{\"periods\": [{\"cashflow\": 1.0, \"interest\": 1.0}]}",
	)
	server.ServeHTTP(responseWriter, request)

	got := responseWriter.Body.String()
	expected := "{\"presentValue\":0.5}"
	assertResponseBody(t, got, expected)
}
func assertStatus(t testing.TB, responseWriter *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := responseWriter.Code

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}

func assertResponseBody(t testing.TB, got, expected string) {
	t.Helper()
	if strings.Trim(got, " \n") != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}
