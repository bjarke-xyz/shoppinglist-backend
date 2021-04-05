package model

type IdentityUser struct {
	ID               string `json:"id"`
	EmailVerified    bool   `json:"emailVerified"`
	Enabled          bool   `json:"enabled"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
	Username         string `json:"username"`
}
