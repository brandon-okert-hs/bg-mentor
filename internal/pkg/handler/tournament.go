package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bcokert/bg-mentor/internal/pkg/database"
	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/model"
	"github.com/bcokert/bg-mentor/pkg/request"
)

type TournamentHandler struct {
	DB           *database.Database
	EntryHandler http.Handler
}

func (h *TournamentHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
				logger.Errorw("Failed to load tournaments", "error", err)
				RespondJSON("Tournament", http.StatusInternalServerError, requesterror.InternalError("Tournament", "An error occurred getting tournaments", req), res)
				return
			}
			RespondJSON("Tournament", http.StatusOK, output, res)
			return
		case http.MethodPost:
			h.POST(res, req)
			return
		default:
			RespondJSON("Tournament", http.StatusNotFound, requesterror.MethodNotFound("Tournament", head, req), res)
			return
		}
	default:
		tid, err := strconv.Atoi(head)
		if err != nil {
			logger.Infow("Failed trying to load tournament", "error", err, "tournamentId", tid)
			RespondJSON("Tournament", http.StatusBadRequest, requesterror.BadRequest("Tournament", "Tournament id must be an integer", req), res)
			return
		}
		output, err = h.GETX(tid)
		if err != nil {
			logger.Errorw("Failed to load tournament", "error", err, "tournamentId", tid)
			RespondJSON("Tournament", http.StatusInternalServerError, requesterror.InternalError("Tournament", "An error occurred getting that tournament", req), res)
			return
		}
		RespondJSON("Tournament", http.StatusOK, output, res)
		return
	}
}

func (h *TournamentHandler) GET() ([]byte, error) {
	tournaments, err := h.DB.GetTournaments()
	if err != nil {
		return nil, err
	}

	return json.Marshal(tournaments)
}

func (h *TournamentHandler) GETX(tournamentID int) ([]byte, error) {
	tournament, err := h.DB.GetTournament(tournamentID)
	if err != nil {
		return nil, err
	}

	return json.Marshal(tournament)
}

func (h *TournamentHandler) POST(w http.ResponseWriter, r *http.Request) {
	email, err := GetMemberEmailFromContext(r.Context())
	if err != nil {
		logger.Errorw("Error getting member token when creating tournament", "error", err)
		RespondJSON("Tournament", http.StatusInternalServerError, requesterror.InternalError("Tournament", "Could not create that tournament", r), w)
		return
	}

	var member *model.Member
	member, err = h.DB.GetMember(email)
	if err != nil {
		logger.Errorw("Error loading member when creating tournament", "error", err)
		RespondJSON("Tournament", http.StatusInternalServerError, requesterror.InternalError("Tournament", "Could not create that tournament", r), w)
		return
	}

	request := model.TournamentPost{}
	if err = fromJSON(r.Body, &request); err != nil {
		RespondJSON("Tournament", http.StatusBadRequest, requesterror.BadRequest("Tournament", err.Error(), r), w)
		return
	}

	if err := request.Validate(); err != nil {
		RespondJSON("Tournament", http.StatusBadRequest, requesterror.BadRequest("Tournament", err.Error(), r), w)
		return
	}

	tournament, err := h.DB.CreateTournament(member.ID, &request)
	if err != nil {
		logger.Errorw("Error creating tournament", "error", err)
		RespondJSON("Tournament", http.StatusInternalServerError, requesterror.InternalError("Tournament", "Could not create that tournament", r), w)
		return
	}

	output, err := json.Marshal(tournament)
	if err != nil {
		logger.Errorw("Error marshalling created tournament", "error", err)
		RespondJSON("Tournament", http.StatusInternalServerError, requesterror.InternalError("Tournament", "Could not retrieve recently created tournament", r), w)
		return
	}
	RespondJSON("Tournament", http.StatusOK, output, w)
}
