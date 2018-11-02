package handler

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
)

type StaticHandler struct {
	StaticFileRoot string
	FileServer     http.Handler
}

func (handler *StaticHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	var err error
	var f *os.File
	var info os.FileInfo

	// Check that the file exists
	f, err = os.Open(path.Clean(fmt.Sprintf("%s/%s", handler.StaticFileRoot, req.URL.Path)))
	if err != nil {
		if os.IsNotExist(err) {
			RespondJSON("Static", http.StatusNotFound, requesterror.FileNotFound("Static", req), res)
		} else {
			RespondJSON("Static", http.StatusInternalServerError, requesterror.FileError("Static", req), res)
		}
		return
	}
	defer f.Close()

	// Check that the file is not a directory
	info, err = f.Stat()
	if err != nil {
		RespondJSON("Static", http.StatusInternalServerError, requesterror.FileError("Static", req), res)
		return
	}
	if info.IsDir() {
		RespondJSON("Static", http.StatusNotFound, requesterror.FileNotFound("Static", req), res)
	} else {
		handler.FileServer.ServeHTTP(res, req)
	}
}
