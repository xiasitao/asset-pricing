package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	a "xiasitao.de/asset-pricing/lib/api"
	i "xiasitao.de/asset-pricing/lib/investments"
)

func requestResponseWriterFactory(path string, body string) (request *http.Request, responseWriter *httptest.ResponseRecorder) {
	request, _ = http.NewRequest(http.MethodPost, path, strings.NewReader(body))
	responseWriter = httptest.NewRecorder()
	return
}

func mockCalcuateGeneralPresentValue(...i.Period) i.CurrencyUnit {
	return 123.5
}

const expectedResponseFromMockCalculateGeneratePresentValue = "{\"presentValue\":123.5}"

func TestPresentValue(t *testing.T) {
	t.Run("test correct request", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory(
			"/general-present-value",
			"{\"periods\": [{\"cashflow\": 1.0, \"interest\": 1.0}]}",
		)
		server := &a.PresentValueServer{CalculateGeneralPresentValue: mockCalcuateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)

		got := responseWriter.Body.String()
		expected := expectedResponseFromMockCalculateGeneratePresentValue

		if strings.Trim(got, " \n") != expected {
			t.Errorf("got %s, expected %s", got, expected)
		}
	})

	t.Run("test malformed request body", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory("/general-present-value", "{malformed}")
		server := &a.PresentValueServer{CalculateGeneralPresentValue: mockCalcuateGeneralPresentValue}
		server.ServeHTTP(responseWriter, request)

		got := responseWriter.Code
		expected := 422

		if got != expected {
			t.Errorf("got %d, expected %d", got, expected)
		}
	})

}
