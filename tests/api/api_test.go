package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xiasitao.de/asset-pricing/lib/api"
)

type SpyInvestmentsServer struct {
	callCount int
}

func (s *SpyInvestmentsServer) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	s.callCount++
}

func TestServer(t *testing.T) {
	mockRequest := httptest.NewRequest(http.MethodPost, "/general-present-value", nil)
	mockResponse := httptest.NewRecorder()

	mockInvestmentsServer := SpyInvestmentsServer{}
	server := api.NewAssetPricingServer(&mockInvestmentsServer)
	server.ServeHTTP(mockResponse, mockRequest)

	if mockInvestmentsServer.callCount != 1 {
		t.Errorf("expected call, didn't receive one")
	}
}
