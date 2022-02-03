package notifications

import "time"

type Message struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	User  User      `json:"user"`
	Time  time.Time `json:"time"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"` //nolint:tagliatelle
	LastName  string `json:"last_name"`  //nolint:tagliatelle
	Email     string `json:"email"`
}
