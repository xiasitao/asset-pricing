package main

import (
	"net/http"

	"xiasitao.de/asset-pricing/lib/api"
)

func main() {
	http.ListenAndServe(":8000", http.HandlerFunc(api.GeneralPresentValueHandler))
}
