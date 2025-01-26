package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	a "xiasitao.de/asset-pricing/lib/api"
)

func requestResponseWriterFactory() (request *http.Request, responseWriter *httptest.ResponseRecorder) {
	request, _ = http.NewRequest(http.MethodGet, "/", nil)
	responseWriter = httptest.NewRecorder()
	return
}

func TestPresentValue(t *testing.T) {
	t.Run("test response code", func(t *testing.T) {
		request, responseWriter := requestResponseWriterFactory()
		a.GeneralPresentValueHandler(responseWriter, request)

		got := responseWriter.Code
		expected := 200

		if got != expected {
			t.Errorf("got %d, expected %d", got, expected)
		}
	})

}
