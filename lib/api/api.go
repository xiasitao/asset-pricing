package api

import (
	"encoding/json"
	"net/http"
)

const UrlPrefix = "/api/asset-pricing"

func ReadRequestBody(body *GeneralPresentValueRequestBody, responseWriter http.ResponseWriter, request *http.Request) error {
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		responseWriter.WriteHeader(422)
		responseWriter.Write([]byte(err.Error()))
	}
	return err
}

func WriteResponseBody(responseWriter http.ResponseWriter, value any) {
	json.NewEncoder(responseWriter).Encode(value)
}
