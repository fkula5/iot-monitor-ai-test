package models

type Device struct {
	ID      string  `json:"id" gorm:"primaryKey"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
	Battery int     `json:"battery"`
	Uptime  string  `json:"uptime"`
	Unit    string  `json:"unit"`
}
