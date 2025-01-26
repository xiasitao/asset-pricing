package main

import (
	"net/http"
	"os"

	"xiasitao.de/asset-pricing/lib/api"
)

func main() {
	router := http.NewServeMux()
	router.Handle(api.UrlPrefix+"/", &api.ProductionInvestmentsServer)
	http.ListenAndServe(getListenAddress(), router)
}

func getListenAddress() string {
	listenAddress := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddress = ":" + val
	}
	return listenAddress
}
