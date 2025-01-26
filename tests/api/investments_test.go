package api

import (
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

func mockCalculateGeneralPresentValue(...i.Period) i.CurrencyUnit {
	return 123.5
}

const expectedResponseFromMockCalculateGeneratePresentValue = "{\"presentValue\":123.5}"

func TestInvestmentsServer(t *testing.T) {
	t.Run("test unknown path", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/unknown-path", "")
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusNotFound)
	})

	t.Run("test wrong method", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		responseWriter := httptest.NewRecorder()
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusMethodNotAllowed)
	})
}

func TestPresentValue(t *testing.T) {
	t.Run("test correct request", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory(
			"/general-present-value",
			"{\"periods\": [{\"cashflow\": 1.0, \"interest\": 1.0}]}",
		)
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)

		got := responseWriter.Body.String()
		expected := expectedResponseFromMockCalculateGeneratePresentValue

		if strings.Trim(got, " \n") != expected {
			t.Errorf("got %s, expected %s", got, expected)
		}
	})

	t.Run("test malformed request body", func(t *testing.T) {
		request, responseWriter := postRequestResponseWriterFactory("/general-present-value", "{malformed}")
		server := &a.InvestmentsServer{CalculateGeneralPresentValue: mockCalculateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)
		assertStatus(t, responseWriter, http.StatusUnprocessableEntity)

	})

}

func assertStatus(t testing.TB, responseWriter *httptest.ResponseRecorder, expected int) {
	t.Helper()
	got := responseWriter.Code

	if got != expected {
		t.Errorf("got %d, expected %d", got, expected)
	}
}
