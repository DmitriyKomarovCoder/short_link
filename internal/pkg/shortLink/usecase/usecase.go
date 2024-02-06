package usecase

import (
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/linkGenerator"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/repository"
	"time"
)

type Usecase struct {
	urlRepo repository.Repository
	log     logger.Logger
	linkGen linkGenerator.LinkHash
}

func NewUsecase(cr repository.Repository, log logger.Logger, hash linkGenerator.LinkHash) *Usecase {
	return &Usecase{
		urlRepo: cr,
		log:     log,
		linkGen: hash,
	}
}

const (
	prefix = "http://"
)

func (u *Usecase) GetUrl(url string) (string, error) {
	shortLink, err := u.urlRepo.GetUrl(url)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Duration(24) * time.Hour)
	if err = u.urlRepo.UpdateTime(expirationTime, shortLink); err != nil {
		return "", err
	}

	return prefix + shortLink, nil
}

func (u *Usecase) CreateLink(longUrl string) (string, error) {
	hashLink := u.linkGen.GenLink(longUrl)
	existsShort, err := u.urlRepo.UrlExistsShort(hashLink)
	if err != nil {
		return "", err
	}

	if existsShort == longUrl {
		return prefix + hashLink, nil
	}
	expirationTime := time.Now().Add(time.Duration(24) * time.Hour)
	if err := u.urlRepo.SaveUrl(longUrl, hashLink, expirationTime); err != nil {
		return "", err
	}

	return prefix + hashLink, err
}
