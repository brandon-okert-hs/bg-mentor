package requesterror

import (
	"fmt"
	"net/http"
	"strings"
)

// PathNotFound should be used when a handler cannot handle the next head in the path.
// It returns an error json response, or an error if something goes wrong
// This is used to handle errors, so we don't want it to fail, thus we manually marshal the json
func PathNotFound(handler, head string, req *http.Request) []byte {
	return []byte(strings.TrimSpace(fmt.Sprintf(`
{
	"message": "The '%s' handler does not have a handler for '%s'",
	"handler": "%s",
	"head":    "%s",
	"method":  "%s",
	"tail":    "%s",
}
	`, handler, head, handler, head, req.Method, req.URL.Path)))
}

// FileNotFound should be used when a handler cannot the given path as a file.
// It returns an error json response, or an error if something goes wrong
// This is used to handle errors, so we don't want it to fail, thus we manually marshal the json
func FileNotFound(handler string, req *http.Request) []byte {
	return []byte(strings.TrimSpace(fmt.Sprintf(`
{
	"message": "The '%s' handler did not find the file '%s'",
	"handler": "%s",
	"method":  "%s",
	"tail":    "%s",
}
	`, handler, req.URL.Path, handler, req.Method, req.URL.Path)))
}

// MethodNotFound should be used when a handler handles the path, but not for the current http method.
// It returns an error json response, or an error if something goes wrong
// This is used to handle errors, so we don't want it to fail, thus we manually marshal the json
func MethodNotFound(handler, head string, req *http.Request) []byte {
	return []byte(strings.TrimSpace(fmt.Sprintf(`
{
	"message": "The '%s' handlers '%s' path does not have a handler for '%s'",
	"handler": "%s",
	"head":    "%s",
	"method":  "%s",
	"tail":    "%s",
}
	`, handler, head, req.Method, handler, head, req.Method, req.URL.Path)))
}
