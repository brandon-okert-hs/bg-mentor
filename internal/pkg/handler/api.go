package handler

import (
	"net/http"

	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/request"
)

type APIHandler struct {
	MemberHandler     http.Handler
	UnitHandler       http.Handler
	TournamentHandler http.Handler
	DABEntryHandler   http.Handler
}

func (h *APIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	switch head {
	case "member":
		h.MemberHandler.ServeHTTP(res, req)
	case "unit":
		h.UnitHandler.ServeHTTP(res, req)
	case "tournament":
		h.TournamentHandler.ServeHTTP(res, req)
	case "dabEntry":
		h.DABEntryHandler.ServeHTTP(res, req)
	default:
		RespondJSON("Api", http.StatusNotFound, requesterror.PathNotFound("Api", head, req), res)
		return
	}
}

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
