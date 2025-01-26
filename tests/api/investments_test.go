package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	a "xiasitao.de/asset-pricing/lib/api"
)

func requestResponseWriterFactory(requestBody string) (request *http.Request, responseWriter *httptest.ResponseRecorder) {
	request, _ = http.NewRequest(http.MethodGet, "/", strings.NewReader(requestBody))
	responseWriter = httptest.NewRecorder()
	return
}

func TestPresentValue(t *testing.T) {
	t.Run("test response code", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory("")
		a.GeneralPresentValueHandler(responseWriter, request)

		got := responseWriter.Code
		expected := 200

		if got != expected {
			t.Errorf("got %d, expected %d", got, expected)
		}
	})

	t.Run("test response body", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory("{\"periods\": [{\"cashflow\": 3000.0}]}")
		a.GeneralPresentValueHandler(responseWriter, request)

		got := responseWriter.Body.String()
		expected := "{\"presentValue\":3000.0}"

		if got != expected {
			t.Errorf("got %s, expected %s", got, expected)
		}
	})

}
