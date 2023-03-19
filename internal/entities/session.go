package entities

import "time"

type Session struct {
	SessionID      string    `db:"session_id"`
	Login          string    `db:"login"`
	CreateDate     time.Time `db:"create_date"`
	ExpireDate     time.Time `db:"expire_date"`
	LastAccessDate time.Time `db:"last_access_date"`
}
