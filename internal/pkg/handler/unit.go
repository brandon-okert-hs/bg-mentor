package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bcokert/bg-mentor/internal/pkg/database"
	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/request"
)

type UnitHandler struct {
	DB *database.Database
}

func (h *UnitHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	var output []byte
	var err error

	switch head {
	case "":
		switch req.Method {
		case http.MethodGet:
			output, err = h.GET()
			if err != nil {
				logger.Errorw("Failed to load units/", "error", err)
				RespondJSON("Unit", http.StatusInternalServerError, requesterror.InternalError("Unit", "An error occurred loading units", req), res)
				return
			}
			RespondJSON("Unit", http.StatusOK, output, res)
			return
		default:
			RespondJSON("Unit", http.StatusNotFound, requesterror.MethodNotFound("Status", head, req), res)
			return
		}
	default:
		RespondJSON("Unit", http.StatusNotFound, requesterror.PathNotFound("Status", head, req), res)
		return
	}
}

func (h *UnitHandler) GET() ([]byte, error) {
	units, err := h.DB.GetUnits()
	if err != nil {
		return nil, err
	}

	return json.Marshal(units)
}
