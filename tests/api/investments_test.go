package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	a "xiasitao.de/asset-pricing/lib/api"
)

func requestResponseWriterFactory(requestBody string) (request *http.Request, responseWriter *httptest.ResponseRecorder) {
	request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
	responseWriter = httptest.NewRecorder()
	return
}

func TestPresentValue(t *testing.T) {
	t.Run("test response code", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory("")
		a.GeneralPresentValueHandler(responseWriter, request)

		got := responseWriter.Code
		expected := 422

		if got != expected {
			t.Errorf("got %d, expected %d", got, expected)
		}
	})

	t.Run("test response body", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory(
			"{\"periods\": [{\"cashflow\": 1.0, \"interest\": 1.0}]}",
		)
		a.GeneralPresentValueHandler(responseWriter, request)

		got := responseWriter.Body.String()
		expected := "{\"presentValue\":0.5}"

		if strings.Trim(got, " \n") != expected {
			t.Errorf("got %s, expected %s", got, expected)
		}
	})

}
