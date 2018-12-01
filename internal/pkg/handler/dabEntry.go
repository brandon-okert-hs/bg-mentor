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

type DABEntryHandler struct {
	DB *database.Database
}

func (h *DABEntryHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	switch head {
	case "":
		switch req.Method {
		case http.MethodGet:
			h.GET(res, req)
			return
		case http.MethodPost:
			h.POST(res, req)
			return
		default:
			RespondJSON("DABEntry", http.StatusNotFound, requesterror.MethodNotFound("DABEntry", head, req), res)
			return
		}
	default:
		switch req.Method {
		case http.MethodPut:
			h.PUT(head, res, req)
			return
		default:
			RespondJSON("DABEntry", http.StatusNotFound, requesterror.MethodNotFound("DABEntry", head, req), res)
			return
		}
	}
}

func (h *DABEntryHandler) GET(w http.ResponseWriter, r *http.Request) {
	tournamentIDRaw := r.URL.Query().Get("tournamentID")
	tournamentID, err := strconv.Atoi(tournamentIDRaw)
	if err != nil {
		logger.Infow("Received illegal tournament id", "error", err, "tournamentIDRaw", tournamentIDRaw)
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", "query param tournamentID is required and must be an int", r), w)
		return
	}

	entries, err := h.DB.GetDABEntries(tournamentID)
	if err != nil {
		logger.Errorw("Error getting DAB entries", "error", err, "tournamentID", tournamentID, "tournamentIDRaw", tournamentIDRaw)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not get dab entries for that tournament", r), w)
		return
	}

	output, err := json.Marshal(entries)
	if err != nil {
		logger.Errorw("Error marshalling entries", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not return dab entries for that tournament", r), w)
		return
	}
	RespondJSON("DABEntry", http.StatusOK, output, w)
}

func (h *DABEntryHandler) POST(w http.ResponseWriter, r *http.Request) {
	email, err := GetMemberEmailFromContext(r.Context())
	if err != nil {
		logger.Errorw("Error getting member token when creating dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not create that dab entry", r), w)
		return
	}

	var member *model.Member
	member, err = h.DB.GetMember(email)
	if err != nil {
		logger.Errorw("Error loading member when creating dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not create that dab entry", r), w)
		return
	}

	request := model.DABEntryPost{}
	if err = fromJSON(r.Body, &request); err != nil {
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", err.Error(), r), w)
		return
	}

	if err := request.Validate(); err != nil {
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", err.Error(), r), w)
		return
	}

	// Verify they own the tournament
	tournament, err := h.DB.GetTournament(request.TournamentID)
	if err != nil {
		logger.Errorw("Error getting tournament to verify member owns it to add a dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not create that dab entry", r), w)
		return
	}
	if tournament.Creator != member.ID {
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", "Only the owner of a tournament can create new entries for it", r), w)
		return
	}

	// Create the entry
	entry, err := h.DB.CreateDABEntry(tournament.ID, &request)
	if err != nil {
		logger.Errorw("Error creating dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not create that dab entry", r), w)
		return
	}

	output, err := json.Marshal(entry)
	if err != nil {
		logger.Errorw("Error marshalling created dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not retrieve recently created dab entry", r), w)
		return
	}
	RespondJSON("DABEntry", http.StatusOK, output, w)
}

func (h *DABEntryHandler) PUT(entryIdRaw string, w http.ResponseWriter, r *http.Request) {
	entryID, err := strconv.Atoi(entryIdRaw)
	if err != nil {
		logger.Infow("Received illegal entry id when editing an entry", "error", err, "entryIdRaw", entryIdRaw)
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", "entryId is required and must be a valid entry id", r), w)
		return
	}

	email, err := GetMemberEmailFromContext(r.Context())
	if err != nil {
		logger.Errorw("Error getting member token when editing dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not edit that dab entry", r), w)
		return
	}

	var member *model.Member
	member, err = h.DB.GetMember(email)
	if err != nil {
		logger.Errorw("Error loading member when editing dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not edit that dab entry", r), w)
		return
	}

	request := model.DABEntryPut{}
	if err = fromJSON(r.Body, &request); err != nil {
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", err.Error(), r), w)
		return
	}

	logger.Infow("reQQQQ", "request", request)

	if err := request.Validate(); err != nil {
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", err.Error(), r), w)
		return
	}

	entry, err := h.DB.GetDABEntry(entryID)
	if err != nil {
		logger.Errorw("Error getting dab entry to verify member can modify it", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not edit that dab entry", r), w)
		return
	}

	tournament, err := h.DB.GetTournament(entry.TournamentID)
	if err != nil {
		logger.Errorw("Error getting tournament to verify member can modify one of its entries", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not edit that dab entry", r), w)
		return
	}

	isEditAllowed := false

	// Tournament owner can always edit
	if tournament.Creator == member.ID {
		isEditAllowed = true
	}

	// Members can only edit if the entry isn't locked by the owner
	if !entry.IsLocked {
		// Members can edit if they are a participant and are editing their own fields
		if entry.Config.Member1 == member.ID {
			if request.Config.Member2 == entry.Config.Member2 &&
				request.Config.Member1 == entry.Config.Member1 && // can't change participants
				request.Config.Member2Race == entry.Config.Member2Race &&
				request.Config.Member2NumBans == entry.Config.Member2NumBans &&
				request.Config.Member1NumBans == entry.Config.Member1NumBans && // can't change num bans
				request.IsLocked == entry.IsLocked && // only owner can lock entries
				request.Config.Member2Ban1 == entry.Config.Member2Ban1 &&
				request.Config.Member2Ban2 == entry.Config.Member2Ban2 &&
				request.Config.Member2Ban3 == entry.Config.Member2Ban3 &&
				request.Config.Member2Ban4 == entry.Config.Member2Ban4 &&
				request.Config.Member2Ban5 == entry.Config.Member2Ban5 &&
				request.Config.Member2Ban6 == entry.Config.Member2Ban6 &&
				request.Config.Member2Confirmed == entry.Config.Member2Confirmed {
				isEditAllowed = true
			}
		}

		if entry.Config.Member2 == member.ID {
			if request.Config.Member1 == entry.Config.Member1 &&
				request.Config.Member2 == entry.Config.Member2 && // can't change participants
				request.Config.Member1Race == entry.Config.Member1Race &&
				request.Config.Member1NumBans == entry.Config.Member1NumBans &&
				request.Config.Member2NumBans == entry.Config.Member2NumBans && // can't change num bans
				request.IsLocked == entry.IsLocked && // only owner can lock entries
				request.Config.Member1Ban1 == entry.Config.Member1Ban1 &&
				request.Config.Member1Ban2 == entry.Config.Member1Ban2 &&
				request.Config.Member1Ban3 == entry.Config.Member1Ban3 &&
				request.Config.Member1Ban4 == entry.Config.Member1Ban4 &&
				request.Config.Member1Ban5 == entry.Config.Member1Ban5 &&
				request.Config.Member1Ban6 == entry.Config.Member1Ban6 &&
				request.Config.Member1Confirmed == entry.Config.Member1Confirmed {
				isEditAllowed = true
			}
		}
	}

	if !isEditAllowed {
		RespondJSON("DABEntry", http.StatusBadRequest, requesterror.BadRequest("DABEntry", "You do not have permission to make that change on this dab entry", r), w)
		return
	}

	// update the entry
	newEntry, err := h.DB.UpdateDABEntry(entry.ID, tournament.ID, &request)
	if err != nil {
		logger.Errorw("Error editing dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not edit that dab entry", r), w)
		return
	}

	output, err := json.Marshal(newEntry)
	if err != nil {
		logger.Errorw("Error marshalling updated dab entry", "error", err)
		RespondJSON("DABEntry", http.StatusInternalServerError, requesterror.InternalError("DABEntry", "Could not retrieve recently edited dab entry", r), w)
		return
	}
	RespondJSON("DABEntry", http.StatusOK, output, w)
}
