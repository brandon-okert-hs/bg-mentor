package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/request"
)

type RootHandler struct {
	StaticFileRoot string
	StatusHandler  http.Handler
	StaticHandler  http.Handler
	AuthHandler    http.Handler
	APIHandler     http.Handler
}

var pageHTML = strings.TrimSpace(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>%s</title>
  </head>
  <body>
  <script type="text/javascript" src="/static/js/%s.js?%s"></script></body>
</html>
`)

func (h *RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	b := make([]byte, 32)
	rand.Read(b)
	hash := base64.StdEncoding.EncodeToString(b)

	switch head {
	case "":
		fmt.Fprintf(res, pageHTML, "Born Gosu Gaming", "index", hash)
	case "tournaments":
		fmt.Fprintf(res, pageHTML, "Tournaments", "tournaments", hash)
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
	case "api":
		h.APIHandler.ServeHTTP(res, req)
	default:
		RespondJSON("Root", http.StatusNotFound, requesterror.PathNotFound("Root", head, req), res)
		return
	}
}
