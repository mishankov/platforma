package session

import "time"

type Session struct {
	ID      string    `db:"id"      json:"id"`
	User    string    `db:"user"    json:"user"`
	Created time.Time `db:"created" json:"created"`
	Expires time.Time `db:"expires" json:"expires"`
}

func (s *Session) IsExpired() bool {
	return s.Expires.Before(time.Now())
}
