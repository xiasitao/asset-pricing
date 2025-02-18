package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

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

func ReadRequestBody(body any, responseWriter http.ResponseWriter, request *http.Request) error {
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

func GetListenerPort(alternatives ...string) int {
	for _, candidate := range alternatives {
		port, err := strconv.Atoi(candidate)
		if err == nil {
			return port
		}
	}
	return 8000
}

func GetListenerAddress() string {
	port_alternative_0, _ := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	port_alternative_1, _ := os.LookupEnv("PORT")
	port := GetListenerPort(port_alternative_0, port_alternative_1)
	address := fmt.Sprintf(":%d", port)
	return address
}
