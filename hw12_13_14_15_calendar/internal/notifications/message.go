package notifications

import "time"

type Message struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	User  int       `json:"user"`
	Time  time.Time `json:"time"`
}
