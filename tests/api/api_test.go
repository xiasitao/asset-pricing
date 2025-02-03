package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiasitao.de/asset-pricing/lib/api"
)

type SpyRouter struct {
	callCount int
}

func (s *SpyRouter) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	s.callCount++
}

func TestAssetPricingServer(t *testing.T) {
	mockRequest := httptest.NewRequest(http.MethodPost, "/general-present-value", nil)
	mockResponse := httptest.NewRecorder()

	mockInvestmentsServer := SpyRouter{}
	server := api.NewAssetPricingServer(&mockInvestmentsServer)
	server.ServeHTTP(mockResponse, mockRequest)

	if mockInvestmentsServer.callCount != 1 {
		t.Errorf("expected call, didn't receive one")
	}
}

func TestGetListenerPort(t *testing.T) {
	assertPort := func(t *testing.T, got, want int) {
		t.Helper()
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
	t.Run("no alternatives", func(t *testing.T) {
		got := api.GetListenerPort()
		want := 8000
		assertPort(t, got, want)
	})
	t.Run("with alternatives", func(t *testing.T) {
		got := api.GetListenerPort("", "1234", "5678")
		want := 1234
		assertPort(t, got, want)
	})

}
func TestGetListenerAddress(t *testing.T) {
	got := api.GetListenerAddress()
	want := ":8000"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
