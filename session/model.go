package session

import "time"

type Session struct {
	ID      string    `json:"id" db:"id"`
	User    string    `json:"user" db:"user"`
	Created time.Time `json:"created" db:"created"`
	Expires time.Time `json:"expires" db:"expires"`
}

func (s *Session) IsExpired() bool {
	return s.Expires.Before(time.Now())
}
