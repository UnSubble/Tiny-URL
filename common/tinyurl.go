package common

import "time"

type TinyURL struct {
	OriginalURL  string
	GeneratedURL string
	GeneratedAt  time.Time
}
