package model

type Search struct {
	Source    string `json:"source"`
	Timestamp int64  `json:"timestamp"`
	Renew     string `json:"renew"`
}
