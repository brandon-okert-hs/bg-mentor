package model

type DABConfig struct {
	ID               int    `json:"id"`
	Member1          int    `json:"member1,omitempty"`
	Member1Race      string `json:"member1Race,omitempty"`
	Member1NumBans   int    `json:"member1NumBans,omitempty"`
	Member1Ban1      string `json:"member1Ban1,omitempty"`
	Member1Ban2      string `json:"member1Ban2,omitempty"`
	Member1Ban3      string `json:"member1Ban3,omitempty"`
	Member1Ban4      string `json:"member1Ban4,omitempty"`
	Member1Ban5      string `json:"member1Ban5,omitempty"`
	Member1Ban6      string `json:"member1Ban6,omitempty"`
	Member1Confirmed bool   `json:"member1Confirmed"`
	Member2          int    `json:"member2,omitempty"`
	Member2Race      string `json:"member2Race,omitempty"`
	Member2NumBans   int    `json:"member2NumBans,omitempty"`
	Member2Ban1      string `json:"member2Ban1,omitempty"`
	Member2Ban2      string `json:"member2Ban2,omitempty"`
	Member2Ban3      string `json:"member2Ban3,omitempty"`
	Member2Ban4      string `json:"member2Ban4,omitempty"`
	Member2Ban5      string `json:"member2Ban5,omitempty"`
	Member2Ban6      string `json:"member2Ban6,omitempty"`
	Member2Confirmed bool   `json:"member2Confirmed"`
}

type DABConfigPut struct {
	Member1          int    `json:"member1,omitempty"`
	Member1Race      string `json:"member1Race,omitempty"`
	Member1NumBans   int    `json:"member1NumBans,omitempty"`
	Member1Ban1      string `json:"member1Ban1,omitempty"`
	Member1Ban2      string `json:"member1Ban2,omitempty"`
	Member1Ban3      string `json:"member1Ban3,omitempty"`
	Member1Ban4      string `json:"member1Ban4,omitempty"`
	Member1Ban5      string `json:"member1Ban5,omitempty"`
	Member1Ban6      string `json:"member1Ban6,omitempty"`
	Member1Confirmed bool   `json:"member1Confirmed"`
	Member2          int    `json:"member2,omitempty"`
	Member2Race      string `json:"member2Race,omitempty"`
	Member2NumBans   int    `json:"member2NumBans,omitempty"`
	Member2Ban1      string `json:"member2Ban1,omitempty"`
	Member2Ban2      string `json:"member2Ban2,omitempty"`
	Member2Ban3      string `json:"member2Ban3,omitempty"`
	Member2Ban4      string `json:"member2Ban4,omitempty"`
	Member2Ban5      string `json:"member2Ban5,omitempty"`
	Member2Ban6      string `json:"member2Ban6,omitempty"`
	Member2Confirmed bool   `json:"member2Confirmed"`
}

func (m *DABConfigPut) Validate() error {
	return nil
}
