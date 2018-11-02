package handler

import (
	"net/http"

	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/request"
)

type StatusHandler struct{}

func (h *StatusHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	switch head {
	case "up":
		switch req.Method {
		case http.MethodGet:
			RespondJSON("Status", http.StatusOK, h.GETUp(), res)
			return
		default:
			RespondJSON("Status", http.StatusNotFound, requesterror.MethodNotFound("Status", head, req), res)
			return
		}
	default:
		RespondJSON("Status", http.StatusNotFound, requesterror.PathNotFound("Status", head, req), res)
		return
	}
}

func (h *StatusHandler) GETUp() []byte {
	return []byte(`{"status": "OK"}`)
}
