package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/bcokert/bg-mentor/pkg/model"
)

func populateDABEntryFromRow(rows *sql.Rows, t *model.DABEntry) error {
	optMember1Bans := make([]sql.NullString, 6)
	optMember2Bans := make([]sql.NullString, 6)
	var optMember1Race sql.NullString
	var optMember2Race sql.NullString
	var optMember1 sql.NullInt64
	var optMember2 sql.NullInt64
	var optMember1NumBans sql.NullInt64
	var optMember2NumBans sql.NullInt64

	err := rows.Scan(
		&t.ID,
		&t.TournamentID,
		&t.IsLocked,
		&optMember1,
		&optMember1Race,
		&optMember1NumBans,
		&optMember1Bans[0],
		&optMember1Bans[1],
		&optMember1Bans[2],
		&optMember1Bans[3],
		&optMember1Bans[4],
		&optMember1Bans[5],
		&t.Config.Member1Confirmed,
		&optMember2,
		&optMember2Race,
		&optMember2NumBans,
		&optMember2Bans[0],
		&optMember2Bans[1],
		&optMember2Bans[2],
		&optMember2Bans[3],
		&optMember2Bans[4],
		&optMember2Bans[5],
		&t.Config.Member2Confirmed,
	)
	if err != nil {
		return fmt.Errorf("Failed to scan dab entry data: %s", err)
	}

	if optMember1Race.Valid {
		t.Config.Member1Race = optMember1Race.String
	}
	if optMember2Race.Valid {
		t.Config.Member2Race = optMember2Race.String
	}

	if optMember1.Valid {
		t.Config.Member1 = int(optMember1.Int64)
	}
	if optMember2.Valid {
		t.Config.Member2 = int(optMember2.Int64)
	}
	if optMember1NumBans.Valid {
		t.Config.Member1NumBans = int(optMember1NumBans.Int64)
	}
	if optMember2NumBans.Valid {
		t.Config.Member2NumBans = int(optMember2NumBans.Int64)
	}

	if optMember1Bans[0].Valid {
		t.Config.Member1Ban1 = optMember1Bans[0].String
	}
	if optMember1Bans[1].Valid {
		t.Config.Member1Ban2 = optMember1Bans[1].String
	}
	if optMember1Bans[2].Valid {
		t.Config.Member1Ban3 = optMember1Bans[2].String
	}
	if optMember1Bans[3].Valid {
		t.Config.Member1Ban4 = optMember1Bans[3].String
	}
	if optMember1Bans[4].Valid {
		t.Config.Member1Ban5 = optMember1Bans[4].String
	}
	if optMember1Bans[5].Valid {
		t.Config.Member1Ban6 = optMember1Bans[5].String
	}

	if optMember2Bans[0].Valid {
		t.Config.Member2Ban1 = optMember2Bans[0].String
	}
	if optMember2Bans[1].Valid {
		t.Config.Member2Ban2 = optMember2Bans[1].String
	}
	if optMember2Bans[2].Valid {
		t.Config.Member2Ban3 = optMember2Bans[2].String
	}
	if optMember2Bans[3].Valid {
		t.Config.Member2Ban4 = optMember2Bans[3].String
	}
	if optMember2Bans[4].Valid {
		t.Config.Member2Ban5 = optMember2Bans[4].String
	}
	if optMember2Bans[5].Valid {
		t.Config.Member2Ban6 = optMember2Bans[5].String
	}

	return nil
}

func (db *Database) GetDABEntry(id int) (*model.DABEntry, error) {
	query := strings.Replace(strings.Replace(`SELECT
		tournament_entries.id,
		tournament_entries.tournamentID,
		tournament_entries.isLocked,

		dab_configs.member1,
		dab_configs.member1Race,
		dab_configs.member1NumBans,
		dab_configs.member1Ban1,
		dab_configs.member1Ban2,
		dab_configs.member1Ban3,
		dab_configs.member1Ban4,
		dab_configs.member1Ban5,
		dab_configs.member1Ban6,
		dab_configs.member1Confirmed,

		dab_configs.member2,
		dab_configs.member2Race,
		dab_configs.member2NumBans,
		dab_configs.member2Ban1,
		dab_configs.member2Ban2,
		dab_configs.member2Ban3,
		dab_configs.member2Ban4,
		dab_configs.member2Ban5,
		dab_configs.member2Ban6,
		dab_configs.member2Confirmed

		FROM tournament_entries INNER JOIN dab_configs
		ON tournament_entries.configId = dab_configs.id
		AND tournament_entries.id = ?;`, "\n", " ", 100), "\t", " ", 1000)
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("Database error loading tournament entries: %s", err)
	}
	defer rows.Close()

	t := model.DABEntry{}
	if rows.Next() {
		if err := populateDABEntryFromRow(rows, &t); err != nil {
			return nil, err
		}
		return &t, nil
	}

	return nil, fmt.Errorf("Did not find dab entry with id '%v'", id)
}

func (db *Database) GetDABEntries(tournamentID int) ([]model.DABEntry, error) {
	query := strings.Replace(strings.Replace(`SELECT
		tournament_entries.id,
		tournament_entries.tournamentID,
		tournament_entries.isLocked,

		dab_configs.member1,
		dab_configs.member1Race,
		dab_configs.member1NumBans,
		dab_configs.member1Ban1,
		dab_configs.member1Ban2,
		dab_configs.member1Ban3,
		dab_configs.member1Ban4,
		dab_configs.member1Ban5,
		dab_configs.member1Ban6,
		dab_configs.member1Confirmed,

		dab_configs.member2,
		dab_configs.member2Race,
		dab_configs.member2NumBans,
		dab_configs.member2Ban1,
		dab_configs.member2Ban2,
		dab_configs.member2Ban3,
		dab_configs.member2Ban4,
		dab_configs.member2Ban5,
		dab_configs.member2Ban6,
		dab_configs.member2Confirmed

		FROM tournament_entries INNER JOIN dab_configs
		ON tournament_entries.configId = dab_configs.id
		AND tournament_entries.tournamentID = ?;`, "\n", " ", 100), "\t", " ", 1000)

	rows, err := db.Query(query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("Database error loading tournament entries: %s", err.Error())
	}

	entries := make([]model.DABEntry, 0, 10)
	for rows.Next() {
		entry := model.DABEntry{}
		if err := populateDABEntryFromRow(rows, &entry); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func (db *Database) CreateDABEntry(tournamentID int, post *model.DABEntryPost) (*model.DABEntry, error) {
	if err := post.Validate(); err != nil {
		return nil, err
	}

	query := "INSERT INTO dab_configs values ()"
	res, err := db.Execute(query)
	if err != nil {
		return nil, fmt.Errorf("Database error creating a dab config for an entry: %s", err)
	}

	dabConfigID, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, fmt.Errorf("Database error obtaining id of dab config for entry: %s", err)
	}

	query = "INSERT INTO tournament_entries (tournamentID, configId) values (?, ?)"
	res, err = db.Execute(query, tournamentID, dabConfigID)
	if err != nil {
		return nil, fmt.Errorf("Database error creating an entry: %s", err)
	}

	entryID, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, fmt.Errorf("Database error obtaining id of entry: %s", err)
	}

	return db.GetDABEntry(int(entryID))
}

func (db *Database) UpdateDABEntry(id int, tournamentId int, put *model.DABEntryPut) (*model.DABEntry, error) {
	if err := put.Validate(); err != nil {
		return nil, err
	}

	entry, err := db.GetDABEntry(int(id))
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve dab entry %v before updating it: %s", id, err)
	}

	query := strings.Replace(strings.Replace(`UPDATE dab_configs SET
		member1 = ?,
		member1Race = ?,
		member1NumBans = ?,
		member1Ban1 = ?,
		member1Ban2 = ?,
		member1Ban3 = ?,
		member1Ban4 = ?,
		member1Ban5 = ?,
		member1Ban6 = ?,
		member1Confirmed = ?,

		member2 = ?,
		member2Race = ?,
		member2NumBans = ?,
		member2Ban1 = ?,
		member2Ban2 = ?,
		member2Ban3 = ?,
		member2Ban4 = ?,
		member2Ban5 = ?,
		member2Ban6 = ?,
		member2Confirmed = ?

		WHERE id = ?;`, "\n", " ", 100), "\t", " ", 1000)

	_, err = db.Execute(
		query,
		nullInt(put.Config.Member1),
		nullString(put.Config.Member1Race),
		nullInt(put.Config.Member1NumBans),
		nullString(put.Config.Member1Ban1),
		nullString(put.Config.Member1Ban2),
		nullString(put.Config.Member1Ban3),
		nullString(put.Config.Member1Ban4),
		nullString(put.Config.Member1Ban5),
		nullString(put.Config.Member1Ban6),
		put.Config.Member1Confirmed,
		nullInt(put.Config.Member2),
		nullString(put.Config.Member2Race),
		nullInt(put.Config.Member2NumBans),
		nullString(put.Config.Member2Ban1),
		nullString(put.Config.Member2Ban2),
		nullString(put.Config.Member2Ban3),
		nullString(put.Config.Member2Ban4),
		nullString(put.Config.Member2Ban5),
		nullString(put.Config.Member2Ban6),
		put.Config.Member2Confirmed,
		entry.Config.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("Database error updating a dab entry: %s", err)
	}

	query = `UPDATE tournament_entries SET isLocked = ? where id = ?;`
	_, err = db.Execute(query, put.IsLocked, tournamentId)
	if err != nil {
		return nil, fmt.Errorf("Database error updating a dab entry: %s", err)
	}

	return db.GetDABEntry(id)
}
