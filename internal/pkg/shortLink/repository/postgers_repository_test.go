package repository

import (
	"errors"
	"fmt"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/models"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestPostgreSQLRepository_SaveUrl(t *testing.T) {
	timeLink := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name           string
		longUrl        string
		shortUrl       string
		expirationTime time.Time
		expectQuery    string
		expectArgs     []interface{}
		expectError    error
	}{
		{
			name:           "Valid case",
			longUrl:        "http://example.com",
			shortUrl:       "abc123",
			expirationTime: timeLink,
			expectQuery:    saveShortLink,
			expectError:    nil,
		},
		{
			name:           "Error case",
			longUrl:        "http://example.com",
			shortUrl:       "abc123",
			expirationTime: timeLink,
			expectQuery:    saveShortLink,
			expectError:    errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			repo := &PostgreSQLRepository{
				Pool: mock,
			}

			escapedQuery := regexp.QuoteMeta(tt.expectQuery)

			mock.ExpectExec(escapedQuery).
				WithArgs(tt.longUrl, tt.shortUrl, tt.expirationTime).
				WillReturnResult(pgxmock.NewResult("INSERT", 1)).
				WillReturnError(tt.expectError)

			err := repo.SaveUrl(tt.longUrl, tt.shortUrl, tt.expirationTime)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLRepository_UpdateTime(t *testing.T) {
	tests := []struct {
		name           string
		expirationTime time.Time
		shortUrl       string
		expectError    error
	}{
		{
			name:           "Valid case",
			expirationTime: time.Now().Add(24 * time.Hour),
			shortUrl:       "abc123",
			expectError:    nil,
		},
		{
			name:           "Error case",
			expirationTime: time.Now().Add(24 * time.Hour),
			shortUrl:       "abc123",
			expectError:    errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			repo := &PostgreSQLRepository{
				Pool: mock,
			}

			escapedQuery := regexp.QuoteMeta(updateTime)

			mock.ExpectExec(escapedQuery).
				WithArgs(tt.expirationTime, tt.shortUrl).
				WillReturnResult(pgxmock.NewResult("UPDATE", 1)).
				WillReturnError(tt.expectError)

			err := repo.UpdateTime(tt.expirationTime, tt.shortUrl)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLRepository_GetUrl(t *testing.T) {
	tests := []struct {
		name        string
		shortUrl    string
		rows        *pgxmock.Rows
		expectError error
		errorPgx    error
		result      string
	}{
		{
			name:        "Valid case",
			shortUrl:    "abc123",
			rows:        pgxmock.NewRows([]string{"original"}).AddRow("www.youtube.com"),
			expectError: nil,
			result:      "www.youtube.com",
		},
		{
			name:        "No such link",
			shortUrl:    "nonexistent",
			rows:        pgxmock.NewRows([]string{}),
			errorPgx:    pgx.ErrNoRows,
			expectError: &models.NoSuchLink{Message: "No such url link: no rows in result set"},
		},
		{
			name:        "Error case",
			shortUrl:    "abc123",
			rows:        pgxmock.NewRows([]string{}),
			expectError: errors.New("mock error"),
			errorPgx:    errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			repo := &PostgreSQLRepository{
				Pool: mock,
			}

			escapedQuery := regexp.QuoteMeta(getLongUrl)

			mock.ExpectQuery(escapedQuery).
				WithArgs(tt.shortUrl).
				WillReturnRows(tt.rows).
				WillReturnError(tt.errorPgx)

			result, err := repo.GetUrl(tt.shortUrl)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.result, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLRepository_UrlExistsShort(t *testing.T) {
	tests := []struct {
		name           string
		shortUrl       string
		rows           *pgxmock.Rows
		expectError    error
		expectedResult string
		errorPgx       error
	}{
		{
			name:           "URL exists",
			shortUrl:       "existing",
			rows:           pgxmock.NewRows([]string{"original"}).AddRow("www.example.com"),
			expectError:    nil,
			expectedResult: "www.example.com",
		},
		{
			name:           "No such URL",
			shortUrl:       "nonexistent",
			rows:           pgxmock.NewRows([]string{}),
			expectError:    nil,
			expectedResult: "",
			errorPgx:       pgx.ErrNoRows,
		},
		{
			name:           "Error case",
			shortUrl:       "error",
			rows:           pgxmock.NewRows([]string{}),
			expectError:    errors.New("mock error"),
			errorPgx:       errors.New("mock error"),
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			repo := &PostgreSQLRepository{
				Pool: mock,
			}

			escapedQuery := regexp.QuoteMeta(getLongUrl)

			mock.ExpectQuery(escapedQuery).
				WithArgs(tt.shortUrl).
				WillReturnRows(tt.rows).
				WillReturnError(tt.errorPgx)

			result, err := repo.UrlExistsShort(tt.shortUrl)

			if tt.expectError != nil {
				assert.EqualError(t, err, tt.expectError.Error())
				assert.Equal(t, tt.expectedResult, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgreSQLRepository_Clear(t *testing.T) {
	tests := []struct {
		name             string
		rowsAffected     int64
		expectedErrorMsg error
		errorPgx         error
	}{
		{
			name:         "Success case",
			rowsAffected: 5,
			errorPgx:     nil,
		},
		{
			name:             "Error case",
			rowsAffected:     0,
			errorPgx:         errors.New("mock error"),
			expectedErrorMsg: fmt.Errorf("error deleting outdated records: %v", errors.New("mock error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, _ := pgxmock.NewPool()

			log, _ := logger.NewLogger("test.log")

			repo := &PostgreSQLRepository{
				Pool: mock,
				Log:  *log,
			}

			mock.ExpectExec(regexp.QuoteMeta(deleteOldLink)).
				WillReturnResult(pgxmock.NewResult("DELETE", tt.rowsAffected)).
				WillReturnError(tt.errorPgx)

			err := repo.Clear()

			if tt.errorPgx != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErrorMsg.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
