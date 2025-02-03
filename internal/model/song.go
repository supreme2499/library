package model

import "time"

type Song struct {
	Id          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"release_date,omitempty"`
	Text        string    `json:"text,omitempty"`
	Link        string    `json:"link,omitempty"`
}
