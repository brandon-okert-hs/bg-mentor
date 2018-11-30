package handler

import (
	"fmt"
	"net/http"

	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/request"
)

type RootHandler struct {
	StaticFileRoot string
	StatusHandler  http.Handler
	StaticHandler  http.Handler
	AuthHandler    http.Handler
	MemberHandler  http.Handler
}

func (h *RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	switch head {
	case "":
		http.ServeFile(res, req, fmt.Sprintf("%s/index.html", h.StaticFileRoot))
	case "favicon.ico":
		http.ServeFile(res, req, fmt.Sprintf("%s/favicon.ico", h.StaticFileRoot))
	case "static":
		h.StaticHandler.ServeHTTP(res, req)
		return
	case "status":
		h.StatusHandler.ServeHTTP(res, req)
		return
	case "auth":
		h.AuthHandler.ServeHTTP(res, req)
		return
	case "member":
		h.MemberHandler.ServeHTTP(res, req)
	default:
		RespondJSON("Root", http.StatusNotFound, requesterror.PathNotFound("Root", head, req), res)
		return
	}
}
