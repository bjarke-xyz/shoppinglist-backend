package models

type HTTPError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
