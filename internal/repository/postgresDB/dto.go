package postgresDB

import "time"

type User struct {
	Email        string
	PasswordHash string
	UrlsLeft     int
}

type Link struct {
	ShortUrl  string
	LongUrl   string
	ExpiresAt time.Time
}
