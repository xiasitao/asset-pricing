package main

import (
	"net/http"

	"xiasitao.de/asset-pricing/lib/api"
)

func main() {
	server := api.NewAssetPricingServer(api.ProductionInvestmentsRouter)
	http.ListenAndServe(api.GetListenerAddress(), server)
}
