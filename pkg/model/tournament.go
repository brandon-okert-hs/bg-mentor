package model

import (
	"fmt"
	"time"
)

type Tournament struct {
	ID            int    `json:"id"`
	Type          string `json:"type"`
	Creator       int    `json:"creator"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	ChallongeLink string `json:"challongeLink,omitempty"`
	LogoLink      string `json:"logoLink,omitempty"`
	StartDate     string `json:"startDate"`
	CheckinDate   string `json:"checkinDate"`
}

type TournamentPost struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	ChallongeLink string `json:"challongeLink,omitempty"`
	LogoLink      string `json:"logoLink,omitempty"`
	StartDate     string `json:"startDate"`
	CheckinDate   string `json:"checkinDate"`
}

func (t *TournamentPost) Validate() error {
	if t.Type != "dab" && t.Type != "handicap" {
		return fmt.Errorf("Type must be one of 'dab', 'handicap'")
	}

	if t.Name == "" {
		return fmt.Errorf("Name must not be empty")
	}

	start, err := time.Parse("2006-01-02 15:04:05", t.StartDate)
	if err != nil {
		return fmt.Errorf("StartDate is not legal: %v", err)
	}

	if start.Before(time.Now().Add(time.Hour)) {
		return fmt.Errorf("StartDate must be an hour or more into the future")
	}

	checkin, err := time.Parse("2006-01-02 15:04:05", t.CheckinDate)
	if err != nil {
		return fmt.Errorf("CheckinDate is not legal: %v", err)
	}

	if checkin.After(start) {
		return fmt.Errorf("CheckinDate must be before StartDate")
	}

	return nil
}
