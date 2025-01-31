package main

import (
	"log"
	"net/http"
	"os"

	"xiasitao.de/asset-pricing/lib/api"
)

func main() {
	server := api.NewAssetPricingServer(&api.ProductionInvestmentsRouter)
	http.ListenAndServe(getListenAddress(), server)
}

func getListenAddress() string {
	listenAddress := ":8080"
	if val, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT"); ok {
		listenAddress = ":" + val
	} else if val, ok := os.LookupEnv("PORT"); ok {
		listenAddress = ":" + val
	}
	log.Printf("listening on %s", listenAddress)
	return listenAddress
}
