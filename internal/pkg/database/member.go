package database

import (
	"database/sql"
	"fmt"

	"github.com/bcokert/bg-mentor/pkg/model"
)

func populateMemberFromRow(rows *sql.Rows, member *model.Member) error {
	var optAvatarURL sql.NullString
	err := rows.Scan(&member.ID, &member.Name, &member.Email, &optAvatarURL)
	if err != nil {
		return fmt.Errorf("Failed to scan member data: %s", err)
	}

	if optAvatarURL.Valid {
		member.AvatarURL = optAvatarURL.String
	}

	return nil
}

func (db *Database) GetMember(email string) (*model.Member, error) {
	query := "select id,name,email,avatarUrl from members where email = ?"
	rows, err := db.Query(query, email)
	if err != nil {
		return nil, fmt.Errorf("Database error loading member: %s", err)
	}
	defer rows.Close()

	member := model.Member{}
	if rows.Next() {
		if err := populateMemberFromRow(rows, &member); err != nil {
			return nil, err
		}
		return &member, nil
	}

	return nil, fmt.Errorf("Did not find member with email '%v'", email)
}

func (db *Database) GetMembers() ([]model.Member, error) {
	query := "select id,name,email,avatarUrl from members order by `id` ASC;"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Database error loading members: %s", err.Error())
	}

	members := make([]model.Member, 0, 20)
	for rows.Next() {
		member := model.Member{}
		if err := populateMemberFromRow(rows, &member); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (db *Database) CreateMember(post *model.MemberPost) (*model.Member, error) {
	if err := post.Validate(); err != nil {
		return nil, err
	}

	query := "INSERT INTO members (name, email, avatarUrl) VALUES(?, ?, ?)"
	_, err := db.Execute(query, post.Name, post.Email, nullString(post.AvatarURL))
	if err != nil {
		return nil, fmt.Errorf("Database error creating a member: %s", err)
	}

	return db.GetMember(post.Email)
}

func (db *Database) CreateMemberIfNotExisting(post *model.MemberPost) (*model.Member, error) {
	if err := post.Validate(); err != nil {
		return nil, err
	}

	member, err := db.GetMember(post.Email)
	if err.Error() == fmt.Sprintf("Did not find member with email '%v'", post.Email) {
		return db.CreateMember(post)
	}
	if err != nil {
		return nil, fmt.Errorf("Database error checking if member exists before creating: %s", err)
	}

	return member, nil
}

func (db *Database) UpdateMember(email string, put *model.MemberPut) (*model.Member, error) {
	if err := put.Validate(); err != nil {
		return nil, err
	}

	rows, err := db.Query("select * from members where email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("An error occurred while checking that member %v exists: %s", email, err)
	}
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("An error occurred while checking retrieving member %v: %s", email, err)
		}
		return nil, fmt.Errorf("Cannot update a member that doesn't exist")
	}
	rows.Close()

	query := "UPDATE members set name=?, avatarUrl=? where email = ?"
	_, err = db.Execute(query, put.Name, nullString(put.AvatarURL), email)
	if err != nil {
		return nil, fmt.Errorf("Database error updating a member: %s", err)
	}

	return db.GetMember(email)
}
