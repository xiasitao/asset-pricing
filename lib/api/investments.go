package api

import "net/http"

func GeneralPresentValueHandler(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte("{\"presentValue\":3000.0}"))
}
