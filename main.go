package main

import (
	"net/http"

	"xiasitao.de/asset-pricing/lib/api"
	"xiasitao.de/asset-pricing/lib/investments"
)

func main() {
	presentValueServer := api.PresentValueServer{CalculateGeneralPresentValue: investments.CalculatePresentValue}
	http.ListenAndServe(":8000", &presentValueServer)
}
