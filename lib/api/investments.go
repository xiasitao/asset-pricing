package api

import (
	"encoding/json"
	"net/http"

	i "xiasitao.de/asset-pricing/lib/investments"
)

type GeneralPresentValueRequestBody struct {
	Periods []i.Period `json:"periods"`
}

type GeneralPresentValueResponseBody struct {
	PresentValue i.CurrencyUnit `json:"presentValue"`
}

func GeneralPresentValueHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var body GeneralPresentValueRequestBody
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		responseWriter.WriteHeader(422)
		responseWriter.Write([]byte(err.Error()))
	}

	result := i.CalculatePresentValue(body.Periods...)
	responseBody := GeneralPresentValueResponseBody{result}
	json.NewEncoder(responseWriter).Encode(responseBody)
}
