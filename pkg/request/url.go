package request

import (
	"path"
	"strings"
)

// ShiftURL takes a url, cleans it up, and returns the head and tail of it
// The head typically corresponds to the current handler, and the tail to sub-handlers
// Assigning req.URL.Path to the tail in a handler will allow chaining this shift in sub handlers
// Mostly copied from https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func ShiftURL(url string) (head, tail string) {
	// Flatten extra slashes, remove trailing slash, remove relative paths
	url = path.Clean("/" + url)

	// Find the first non-root slash, so we know where to slice for head and tail
	i := strings.Index(url[1:], "/") + 1
	if i <= 0 {
		return url[1:], "/"
	}
	return url[1:i], url[i:]
}
