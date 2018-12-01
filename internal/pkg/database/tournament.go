package database

import (
	"database/sql"
	"fmt"

	"github.com/bcokert/bg-mentor/pkg/model"
)

func populateTournamentFromRow(rows *sql.Rows, t *model.Tournament) error {
	var optDescription sql.NullString
	var optChallongeLink sql.NullString
	var optLogoLink sql.NullString
	err := rows.Scan(&t.ID, &t.Type, &t.Creator, &t.Name, &optDescription, &optChallongeLink, &optLogoLink, &t.StartDate, &t.CheckinDate)
	if err != nil {
		return fmt.Errorf("Failed to scan member data: %s", err)
	}

	if optDescription.Valid {
		t.Description = optDescription.String
	}
	if optChallongeLink.Valid {
		t.ChallongeLink = optChallongeLink.String
	}
	if optLogoLink.Valid {
		t.LogoLink = optLogoLink.String
	}

	return nil
}

func (db *Database) GetTournament(id int) (*model.Tournament, error) {
	query := "select id,type,creator,name,description,challongeLink,logoLink,startDate,checkinDate from tournaments where id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Database error loading tournament: %s", err)
	}
	defer rows.Close()

	t := model.Tournament{}
	if rows.Next() {
		if err := populateTournamentFromRow(rows, &t); err != nil {
			return nil, err
		}
		return &t, nil
	}

	return nil, fmt.Errorf("Did not find tournament with id '%v'", id)
}

func (db *Database) GetTournaments() ([]model.Tournament, error) {
	query := "select id,type,creator,name,description,challongeLink,logoLink,startDate,checkinDate from tournaments order by `id` ASC;"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Database error loading tournaments: %s", err.Error())
	}

	tournaments := make([]model.Tournament, 0, 20)
	for rows.Next() {
		tournament := model.Tournament{}
		if err := populateTournamentFromRow(rows, &tournament); err != nil {
			return nil, err
		}

		tournaments = append(tournaments, tournament)
	}

	return tournaments, nil
}

func (db *Database) CreateTournament(creatorID int, post *model.TournamentPost) (*model.Tournament, error) {
	if err := post.Validate(); err != nil {
		return nil, err
	}

	query := "INSERT INTO tournaments (type, name, description, challongeLink, logoLink, startDate, checkinDate, creator) VALUES(?, ?, ?, ?, ?, ?, ?, ?)"
	res, err := db.Execute(query, post.Type, post.Name, nullString(post.Description), nullString(post.ChallongeLink), nullString(post.LogoLink), post.StartDate, post.CheckinDate, creatorID)
	if err != nil {
		return nil, fmt.Errorf("Database error creating a tournament: %s", err)
	}

	recentID, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, fmt.Errorf("Database error obtaining id of tournmanet: %s", err)
	}

	return db.GetTournament(int(recentID))
}
