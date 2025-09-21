package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestUrlShortener_GetShortenUrl(t *testing.T) {
	ctx := context.Background()
	mockedLog := zaptest.NewLogger(t)

	mockedGetID := "GetID"
	urlRepo := new(UrlRepositoryMock)
	// Test case 1: URL is found in the database
	argFullUrl1 := "https://fullUrl1.com"
	mockRes1Id := 1
	var mockRes1Err error = nil
	urlRepo.On(mockedGetID, ctx, argFullUrl1).Return(mockRes1Id, mockRes1Err).Once()

	// Test case 2: URL is not found in the database
	argFullUrl2 := "https://fullUrl2.com"
	mockRes2Id := 0
	mockRes2Err := ErrUrlNotFound
	urlRepo.On(mockedGetID, ctx, argFullUrl2).Return(mockRes2Id, mockRes2Err).Once()

	shortUrl, err := encodeID(1)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		fullUrl string
		want    string
		wantErr bool
	}{
		{
			name:    "Test case 1: URL is found in the database",
			fullUrl: argFullUrl1,
			want:    shortUrl,
			wantErr: false,
		},
		{
			name:    "Test case 2: URL is not found in the database",
			fullUrl: argFullUrl2,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlShortener(urlRepo, mockedLog)
			got, gotErr := u.GetShortenUrl(context.Background(), tt.fullUrl)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetShortenUrl() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetShortenUrl() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetShortenUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlShortener_GetFullUrl(t *testing.T) {
	ctx := context.Background()
	mockedLog := zaptest.NewLogger(t)

	mockedGetUrlByID := "GetUrlByID"
	urlRepo := new(UrlRepositoryMock)

	// Test case 1: URL is found in the database
	mockArg1Id := int64(1)
	mockRes1Url := "https://fullUrl1.com"
	var mockRes1Err error = nil
	urlRepo.On(mockedGetUrlByID, ctx, mockArg1Id).Return(mockRes1Url, mockRes1Err).Once()
	arg1ShortenUrl, err := encodeID(mockArg1Id)
	if err != nil {
		t.Fatal(err)
	}

	// Test case 2: URL is not found in the database
	mockArg2Id := int64(2)
	if err != nil {
		t.Fatal(err)
	}
	mockRes2Url := ""
	mockRes2Err := ErrUrlNotFound
	urlRepo.On(mockedGetUrlByID, ctx, mockArg2Id).Return(mockRes2Url, mockRes2Err).Once()
	arg2ShortenUrl, err := encodeID(mockArg2Id)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		shortenUrl string
		want       string
		wantErr    bool
	}{
		{
			name:       "URL is found in the database",
			shortenUrl: arg1ShortenUrl,
			want:       mockRes1Url,
			wantErr:    false,
		},
		{
			name:       "URL is not found in the database",
			shortenUrl: arg2ShortenUrl,
			want:       mockRes2Url,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlShortener(urlRepo, mockedLog)
			got, gotErr := u.GetFullUrl(context.Background(), tt.shortenUrl)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetFullUrl() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetFullUrl() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetFullUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlShortener_CreateUrl(t *testing.T) {
	mockLog := zaptest.NewLogger(t)
	ctx := context.Background()

	mockedGetID := "GetID"
	mockedSaveUrl := "SaveUrl"
	urlRepository := new(UrlRepositoryMock)

	// Test case 1: URL is found in the database
	mockArg1Url := "https://fullUrl1.com"
	mockRes1Id := 1
	var mockRes1Err error = nil
	urlRepository.On(mockedGetID, ctx, mockArg1Url).Return(mockRes1Id, mockRes1Err).Once()
	mockRes1ShortUrl, err := encodeID(int64(mockRes1Id))
	if err != nil {
		t.Fatal(err)
	}
	// Test case 2: URL is not found in the database
	mockArg2Url := "https://fullUrl2.com"
	mockRes2Id := 2
	var mockRes2Err error = nil
	// First call -- URL not found
	urlRepository.On(mockedGetID, ctx, mockArg2Url).Return(0, ErrUrlNotFound).Once()
	// Save URL
	urlRepository.On(mockedSaveUrl, ctx, mockArg2Url).Return(nil).Once()
	// Second call -- URL found
	urlRepository.On(mockedGetID, ctx, mockArg2Url).Return(mockRes2Id, mockRes2Err).Once()
	mockRes2ShortUrl, err := encodeID(int64(mockRes2Id))

	// Test case 3: Error on save
	mockArgTC3Url := "https://fullUrl3.com"
	// First call -- URL not found
	urlRepository.On(mockedGetID, ctx, mockArgTC3Url).Return(0, ErrUrlNotFound).Once()
	// Error on save
	urlRepository.On(mockedSaveUrl, ctx, mockArgTC3Url).Return(errors.New("save error")).Once()
	mockResTC3Res := ""

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		urlRepo UrlRepository
		log     *zap.Logger
		// Named input parameters for target function.
		fullUrl string
		want    string
		err     error
	}{
		// Test case 1: URL is found in the database
		{
			name:    "Test case 1: URL is found in the database",
			urlRepo: urlRepository,
			log:     mockLog,
			fullUrl: mockArg1Url,
			want:    mockRes1ShortUrl,
			err:     mockRes1Err,
		},
		// Test case 2: URL is not found in the database
		{
			name:    "Test case 2: URL is not found in the database",
			urlRepo: urlRepository,
			log:     mockLog,
			fullUrl: mockArg2Url,
			want:    mockRes2ShortUrl,
			err:     mockRes2Err,
		},
		// Test case 3: Error on save
		{
			name:    "Test case 3: Error on save",
			urlRepo: urlRepository,
			log:     mockLog,
			fullUrl: mockArgTC3Url,
			want:    mockResTC3Res,
			err:     errors.New("save error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUrlShortener(tt.urlRepo, tt.log)
			got, err := u.CreateUrl(ctx, tt.fullUrl)
			assert.Equal(t, tt.err, err)
			if got != tt.want {
				t.Errorf("UrlShortener.CreateUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
