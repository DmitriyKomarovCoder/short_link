package repository

import "time"

type Repository interface {
	Connect() error
	Close() error
	Clear() error

	GetUrl(url string) (string, error)
	UrlExistsShort(url string) (string, error)
	SaveUrl(longUrl, shortUrl string, expirationTime time.Time) error
	UpdateTime(expirationTime time.Time, shortLink string) error
}
