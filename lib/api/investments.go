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
	err := readRequestBody(&body, responseWriter, request)
	if err != nil {
		return
	}

	presentValue := i.CalculatePresentValue(body.Periods...)
	json.NewEncoder(responseWriter).Encode(GeneralPresentValueResponseBody{presentValue})
}

func readRequestBody(body *GeneralPresentValueRequestBody, responseWriter http.ResponseWriter, request *http.Request) error {
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		responseWriter.WriteHeader(422)
		responseWriter.Write([]byte(err.Error()))
	}
	return err
}
