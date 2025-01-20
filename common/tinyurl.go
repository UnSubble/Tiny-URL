package common

import "time"

type TinyURL struct {
	OriginalURL string
	ShortCode   string
	GeneratedAt time.Time
}
