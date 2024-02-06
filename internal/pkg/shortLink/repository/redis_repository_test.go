package repository

import (
	"errors"
	"fmt"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSaveUrl(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	log, _ := logger.NewLogger("test.log")
	repo := &RedisRepository{
		Client: mockClient,
		Log:    *log,
	}

	shortUrl := "short"
	longUrl := "https://www.example.com"
	expirationTime := time.Now().Add(1 * time.Hour)

	mock.ExpectSet(shortUrl, longUrl, 0).SetVal(shortUrl)
	mock.ExpectExpire(shortUrl, expirationTime.Sub(time.Now())).SetVal(true)

	err := repo.SaveUrl(longUrl, shortUrl, expirationTime)

	assert.Nil(t, err, "Expected no error, got %v", err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestSaveUrl_SetError(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	log, _ := logger.NewLogger("test.log")
	repo := &RedisRepository{
		Client: mockClient,
		Log:    *log,
	}

	shortUrl := "short"
	longUrl := "https://www.example.com"
	expirationTime := time.Now().Add(1 * time.Hour)

	expectedError := fmt.Errorf("set error")
	mock.ExpectSet(shortUrl, longUrl, 0).SetErr(expectedError)

	err := repo.SaveUrl(longUrl, shortUrl, expirationTime)

	assert.EqualError(t, err, expectedError.Error(), "Expected error does not match actual error")
}

func TestUrlExistsShort_Success(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	shortUrl := "short"
	longUrl := "www.example.com"
	mock.ExpectGet(shortUrl).SetVal(longUrl)

	existsUrl, err := repo.UrlExistsShort(shortUrl)

	assert.Nil(t, err, "Expected no error, got %v", err)
	assert.Equal(t, longUrl, existsUrl, "Expected %s, got %s", longUrl, existsUrl)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestUrlExistsShort_Error(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	shortUrl := "short"
	expectedError := errors.New("something went wrong")
	mock.ExpectGet(shortUrl).SetErr(expectedError)

	existsUrl, err := repo.UrlExistsShort(shortUrl)

	assert.Error(t, err, "Expected an error")

	assert.Equal(t, existsUrl, "", "Expected %s, got %s", existsUrl, "")

	assert.Equal(t, expectedError, err, "Expected %v, got %v", expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestUpdateTime_Success(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	shortLink := "short"
	expirationTime := time.Now().Add(1 * time.Hour)
	duration := expirationTime.Sub(time.Now())

	mock.ExpectExpire(shortLink, duration).SetVal(true)

	err := repo.UpdateTime(expirationTime, shortLink)

	assert.Nil(t, err, "Expected no error, got %v", err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestUpdateTime_Error(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	shortLink := "short"
	expirationTime := time.Now().Add(1 * time.Hour)
	duration := expirationTime.Sub(time.Now())
	expectedError := errors.New("something went wrong")

	mock.ExpectExpire(shortLink, duration).SetErr(expectedError)

	err := repo.UpdateTime(expirationTime, shortLink)

	expectedErrorMessage := errors.New("error updating expiration time for key short: something went wrong")

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, expectedErrorMessage, err, "Expected %v, got %v", expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestGetUrl_Success(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	url := "short"
	longUrl := "https://www.example.com"
	mock.ExpectGet(url).SetVal(longUrl)

	result, err := repo.GetUrl(url)

	assert.Nil(t, err, "Expected no error, got %v", err)
	assert.Equal(t, longUrl, result, "Expected long URL %s, got %s", longUrl, result)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestGetUrl_NotFound(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	url := "short"
	expectedError := redis.Nil
	mock.ExpectGet(url).SetErr(expectedError)

	result, err := repo.GetUrl(url)

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, "", "Expected %s, got %s", result, "")

	assert.IsType(t, &models.NoSuchLink{}, err, "Expected error type to be NoSuchLink")
	assert.Equal(t, fmt.Sprintf("No such url link: %v", expectedError), err.Error(), "Expected error message %q, got %q", fmt.Sprintf("No such url link: %v", expectedError), err.Error())

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}

func TestGetUrl_Error(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()

	repo := &RedisRepository{
		Client: mockClient,
	}

	url := "short"
	expectedError := errors.New("something went wrong")
	mock.ExpectGet(url).SetErr(expectedError)

	result, err := repo.GetUrl(url)

	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, "", "Expected %s, got %s", result, "")
	assert.Equal(t, expectedError, err, "Expected %v, got %v", expectedError, err)

	assert.NoError(t, mock.ExpectationsWereMet(), "Expectations were not met")
}
