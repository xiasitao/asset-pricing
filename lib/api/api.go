package api

import (
	"encoding/json"
	"net/http"
)

func ReadRequestBody(body *GeneralPresentValueRequestBody, responseWriter http.ResponseWriter, request *http.Request) error {
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		responseWriter.WriteHeader(422)
		responseWriter.Write([]byte(err.Error()))
	}
	return err
}
