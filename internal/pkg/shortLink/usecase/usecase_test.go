package usecase

import (
	"errors"
	"github.com/DmitriyKomarovCoder/short_link/common/logger"
	"github.com/DmitriyKomarovCoder/short_link/internal/pkg/linkGenerator"
	mock "github.com/DmitriyKomarovCoder/short_link/internal/pkg/shortLink/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUrl(t *testing.T) {
	testCases := []struct {
		name            string
		url             string
		repoExpectation func(*mock.MockRepository)
		expectedShort   string
		expectedError   error
	}{
		{
			name: "Success",
			url:  "http://example.com",
			repoExpectation: func(repo *mock.MockRepository) {
				repo.EXPECT().GetUrl("http://example.com").Return("abc123", nil)
				repo.EXPECT().UpdateTime(gomock.Any(), "abc123").Return(nil)
			},
			expectedShort: "http://abc123",
			expectedError: nil,
		},
		{
			name: "Repository Error Get",
			url:  "http://example.com",
			repoExpectation: func(repo *mock.MockRepository) {
				repo.EXPECT().GetUrl("http://example.com").Return("", errors.New("repository error"))
			},
			expectedShort: "",
			expectedError: errors.New("repository error"),
		},
		{
			name: "Repository Error Update",
			url:  "http://example.com",
			repoExpectation: func(repo *mock.MockRepository) {
				repo.EXPECT().GetUrl("http://example.com").Return("abc123", nil)
				repo.EXPECT().UpdateTime(gomock.Any(), "abc123").Return(errors.New("repository error"))
			},
			expectedShort: "",
			expectedError: errors.New("repository error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.repoExpectation(mockRepo)

			linkGen := linkGenerator.NewLinkHash("123456789", 10)
			uc := NewUsecase(mockRepo, logger.Logger{}, linkGen)

			shortURL, err := uc.GetUrl(tc.url)

			assert.Equal(t, tc.expectedShort, shortURL)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreateLink(t *testing.T) {
	testCases := []struct {
		name           string
		longUrl        string
		repoExpectFunc func(*mock.MockRepository)
		expectedShort  string
		expectedError  error
	}{
		{
			name:    "Success",
			longUrl: "http://example.com",
			repoExpectFunc: func(repo *mock.MockRepository) {
				repo.EXPECT().UrlExistsShort(gomock.Any()).Return("", nil)
				repo.EXPECT().SaveUrl(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedShort: "http://8545999729",
			expectedError: nil,
		},
		{
			name:    "Repository Error UrlExistsShort",
			longUrl: "http://example.com",
			repoExpectFunc: func(repo *mock.MockRepository) {
				repo.EXPECT().UrlExistsShort(gomock.Any()).Return("", errors.New("repository error"))
			},
			expectedShort: "",
			expectedError: errors.New("repository error"),
		},
		{
			name:    "Repository Error SaveUrl",
			longUrl: "http://example.com",
			repoExpectFunc: func(repo *mock.MockRepository) {
				repo.EXPECT().UrlExistsShort(gomock.Any()).Return("", nil)
				repo.EXPECT().SaveUrl(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
			},
			expectedShort: "",
			expectedError: errors.New("repository error"),
		},
		{
			name:    "Short URL Exists",
			longUrl: "http://example.com",
			repoExpectFunc: func(repo *mock.MockRepository) {
				repo.EXPECT().UrlExistsShort(gomock.Any()).Return("http://example.com", nil)
			},
			expectedShort: "http://8545999729",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockRepository(ctrl)
			tc.repoExpectFunc(mockRepo)

			linkGen := linkGenerator.NewLinkHash("123456789", 9)
			uc := NewUsecase(mockRepo, logger.Logger{}, linkGen)

			shortURL, err := uc.CreateLink(tc.longUrl)

			assert.Equal(t, tc.expectedShort, shortURL)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
