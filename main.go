package main

import (
	"net/http"
	"os"

	"xiasitao.de/asset-pricing/lib/api"
)

func main() {
	http.ListenAndServe(getListenAddress(), &api.ProductionInvestmentsServer)
}

func getListenAddress() string {
	listenAddress := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddress = ":" + val
	}
	return listenAddress
}
