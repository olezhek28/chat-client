package model

import "time"

type Message struct {
	From      string
	Text      string
	CreatedAt time.Time
}
