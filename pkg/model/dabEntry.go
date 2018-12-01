package model

import (
	"fmt"
)

type DABEntry struct {
	ID           int       `json:"id"`
	TournamentID int       `json:"tournamentId"`
	Config       DABConfig `json:"config"`
	IsLocked     bool      `json:"isLocked"`
}

type DABEntryPost struct {
	TournamentID int `json:"tournamentId"`
}

func (m *DABEntryPost) Validate() error {
	if m.TournamentID == 0 {
		return fmt.Errorf("TournamentID must be set")
	}

	return nil
}

type DABEntryPut struct {
	Config   DABConfigPut `json:"config"`
	IsLocked bool         `json:"isLocked"`
}

func (m *DABEntryPut) Validate() error {
	return m.Config.Validate()
}
