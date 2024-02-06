package repository

import (
	"context"
	"fmt"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/models"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisRepository struct {
	Addr     string
	Password string
	DB       int
	Client   *redis.Client
	Log      logger.Logger
	ctx      context.Context
}

func NewRedisRepository(addr string, db int, log logger.Logger) *RedisRepository {
	return &RedisRepository{
		Addr: addr,
		DB:   db,
		Log:  log,
	}
}

func (r *RedisRepository) Connect() error {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	_, err := r.Client.Ping(context.Background()).Result()
	return err
}

func (r *RedisRepository) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

func (r *RedisRepository) Clear() error {
	return nil
}

func (r *RedisRepository) GetUrl(url string) (string, error) {
	longUrl, err := r.Client.Get(context.Background(), url).Result()
	if err != nil {
		if err == redis.Nil {
			return "", &models.NoSuchLink{Message: fmt.Sprintf("No such url link: %v", err)}
		}
		return "", err
	}

	return longUrl, nil
}

func (r *RedisRepository) UrlExistsShort(url string) (string, error) {
	existingURL, err := r.Client.Get(context.Background(), url).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return existingURL, nil
}

func (r *RedisRepository) SaveUrl(longUrl, shortUrl string, expirationTime time.Time) error {
	err := r.Client.Set(context.Background(), shortUrl, longUrl, 0).Err()
	if err != nil {
		return err
	}

	duration := expirationTime.Sub(time.Now())
	err = r.Client.Expire(context.Background(), shortUrl, duration).Err()
	if err != nil {
		return fmt.Errorf("error setting expiration time for switch %s: %v", shortUrl, err)
	}
	return nil
}

func (p *RedisRepository) UpdateTime(expirationTime time.Time, shortLink string) error {
	duration := expirationTime.Sub(time.Now())
	err := p.Client.Expire(context.Background(), shortLink, duration).Err()
	if err != nil {
		return fmt.Errorf("error updating expiration time for key %s: %v", shortLink, err)
	}

	return nil
}
