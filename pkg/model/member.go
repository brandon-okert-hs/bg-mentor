package model

import "fmt"

type Member struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

type MemberPost struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

func (m *MemberPost) Validate() error {
	if m.Name == "" || m.Email == "" {
		return fmt.Errorf("Name and email are both required for new members")
	}

	return nil
}

type MemberPut struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

func (m *MemberPut) Validate() error {
	if m.Name == "" || m.Email == "" {
		return fmt.Errorf("Name and email are both required for new members")
	}

	return nil
}
