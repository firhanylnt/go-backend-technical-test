package models

type Event struct {
	ID          int    `json:"id"`
	Name       	string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
