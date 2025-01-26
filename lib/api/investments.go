package api

import (
	"encoding/json"
	"net/http"
)

type GeneralPresentValueBody struct {
	Periods []Period `json:"periods"`
}

func GeneralPresentValueHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var body GeneralPresentValueBody
	err := json.NewDecoder(request.Body).Decode(body)
	if err != nil {
		responseWriter.WriteHeader(422)
		responseWriter.Write([]byte(err.Error()))
	}

	responseWriter.Write([]byte("{\"presentValue\":3000.0}"))
}
