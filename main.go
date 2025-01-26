package main

import (
	"net/http"
	"os"

	"xiasitao.de/asset-pricing/lib/api"
)

func main() {
	http.Handle(api.UrlPrefix, &api.ProductionInvestmentsServer)
	http.ListenAndServe(getListenAddress(), nil)
}

func getListenAddress() string {
	listenAddress := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddress = ":" + val
	}
	return listenAddress
}
