package database

import (
	"database/sql"
	"fmt"

	"github.com/bcokert/bg-mentor/pkg/model"
)

func populateUnitFromRow(rows *sql.Rows, unit *model.Unit) error {
	err := rows.Scan(&unit.Name, &unit.Race, &unit.DABCanBan)
	if err != nil {
		return fmt.Errorf("Failed to scan unit data: %s", err)
	}

	return nil
}

func (db *Database) GetUnits() ([]model.Unit, error) {
	query := "select name,race,dabCanBan from units order by `race` ASC;"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Database error loading units: %s", err.Error())
	}

	units := make([]model.Unit, 0, 44)
	for rows.Next() {
		unit := model.Unit{}
		if err := populateUnitFromRow(rows, &unit); err != nil {
			return nil, err
		}

		units = append(units, unit)
	}

	return units, nil
}
