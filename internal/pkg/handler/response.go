package handler

import (
	"net/http"
)

func RespondJSON(handlerName string, status int, output []byte, res http.ResponseWriter) {
	res.WriteHeader(status)
	res.Header().Add("Content-Type", "application/json")

	_, err := res.Write(output)
	if err != nil {
		logger.Errorw("Error writing response",
			"handler", handlerName,
			"error", err)
	}
}
