package models

type Bridge struct {
	Destination string `json:"destination"`
	PermWA      bool   `json:"perm_wa"`
	PermTG      bool   `json:"perm_tg"`
}
