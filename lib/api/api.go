package api

import (
	"encoding/json"
	"net/http"
)

const UrlPrefix = ""

type AssetPricingServer struct {
	http.Handler
}

func NewAssetPricingServer(investmentsRouter http.Handler) *AssetPricingServer {
	server := new(AssetPricingServer)
	router := http.NewServeMux()
	router.Handle("/", investmentsRouter)
	server.Handler = router
	return server
}

func ReadRequestBody(body *GeneralPresentValueRequestBody, responseWriter http.ResponseWriter, request *http.Request) error {
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusUnprocessableEntity)
		responseWriter.Write([]byte(err.Error()))
	}
	return err
}

func WriteResponseBody(responseWriter http.ResponseWriter, value any) {
	json.NewEncoder(responseWriter).Encode(value)
}
