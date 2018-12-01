package handler

import (
	"encoding/json"
	"io"
)

// Reads a readCloser, typically the body of a request, into the given struct
// The address of the model should be passed, not the value itself
func fromJSON(body io.ReadCloser, obj interface{}) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
