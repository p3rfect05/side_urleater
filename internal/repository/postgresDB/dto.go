package postgresDB

import "time"

type User struct {
	Email        string
	PasswordHash string
	UrlsLeft     string
}

type Link struct {
	ShortUrl  string
	LongUrl   string
	ExpiresAt time.Time
}
