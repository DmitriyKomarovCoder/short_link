package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

const (
	getLongUrl    = "SELECT original FROM urls WHERE short = $1;"
	saveShortLink = "INSERT INTO urls (original, short, expiration_time) VALUES ($1, $2, $3);"
	updateTime    = "UPDATE urls SET expiration_time = $1 WHERE short = $2;"
	deleteOldLink = "DELETE FROM urls WHERE expiration_time < NOW();"
)

type PostgreSQLRepository struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Pool     Querier
	Log      logger.Logger
	Ctx      context.Context
}

func NewPostgreSQLRepository(ctx context.Context, host string, port int, username, password, dbName string, log logger.Logger) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		DBName:   dbName,
		Log:      log,
		Ctx:      ctx,
	}
}

func (p *PostgreSQLRepository) Connect() error {
	fmt.Println(p.DBName)
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.Username, p.Password, p.DBName)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}

	p.Pool, err = pgxpool.ConnectConfig(p.Ctx, poolConfig)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQLRepository) Close() error {
	if p.Pool != nil {
		p.Pool.Close()
	}
	return nil
}

func (r *PostgreSQLRepository) Clear() error {
	result, err := r.Pool.Exec(context.Background(), deleteOldLink)
	if err != nil {
		return fmt.Errorf("error deleting outdated records: %v", err)
	}

	numDeleted := result.RowsAffected()

	r.Log.Infof("Outdated entries have been successfully deleted, it has been deleted: %v", numDeleted)
	return nil
}

func (p *PostgreSQLRepository) SaveUrl(longUrl, shortUrl string, expirationTime time.Time) error {
	_, err := p.Pool.Exec(context.Background(), saveShortLink, longUrl, shortUrl, expirationTime)

	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQLRepository) UpdateTime(expirationTime time.Time, shortLink string) error {
	_, err := p.Pool.Exec(context.Background(), updateTime, expirationTime, shortLink)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgreSQLRepository) GetUrl(url string) (string, error) {
	var existingURL string
	err := p.Pool.QueryRow(context.Background(), getLongUrl, url).Scan(&existingURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", &models.NoSuchLink{Message: fmt.Sprintf("No such url link: %v", err)}
		}
		return "", err
	}

	return existingURL, nil
}

// проверка если такой хеш существует
func (p *PostgreSQLRepository) UrlExistsShort(url string) (string, error) {
	var existingURL string
	err := p.Pool.QueryRow(context.Background(), getLongUrl, url).Scan(&existingURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return existingURL, nil
}
