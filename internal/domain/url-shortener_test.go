package domain

import (
	"context"
	"testing"

	"go.uber.org/zap/zaptest"
)

func (u *UrlRepositoryMock) GetID(ctx context.Context, fullUrl string) (id int64, err error) {
	args := u.Called(ctx, fullUrl)
	return int64(args.Int(0)), args.Error(1)
}

func (u *UrlRepositoryMock) GetUrlByID(ctx context.Context, id int64) (fullUrl string, err error) {
	args := u.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (u *UrlRepositoryMock) SaveUrl(ctx context.Context, fullUrl string) (err error) {
	args := u.Called(ctx, fullUrl)
	return args.Error(0)
}

func TestUrlShortener_GetShortenUrl(t *testing.T) {
	ctx := context.Background()
	mockedLog := zaptest.NewLogger(t)

	mockedGetID := "GetID"
	urlRepo := new(UrlRepositoryMock)
	// Test case 1: URL is found in the database
	argFullUrlTC1 := "https://fullUrl1.com"
	mockResTC1Id := 1
	var mockResTC1Err error = nil
	urlRepo.On(mockedGetID, ctx, argFullUrlTC1).Return(mockResTC1Id, mockResTC1Err).Once()

	// Test case 2: URL is not found in the database
	argFullUrlTC2 := "https://fullUrl2.com"
	mockResTC2Id := 0
	mockResTC2Err := ErrUrlNotFound
	urlRepo.On(mockedGetID, ctx, argFullUrlTC2).Return(mockResTC2Id, mockResTC2Err).Once()

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
			fullUrl: argFullUrlTC1,
			want:    shortUrl,
			wantErr: false,
		},
		{
			name:    "Test case 2: URL is not found in the database",
			fullUrl: argFullUrlTC2,
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
	mockArgTC1Id := int64(1)
	mockResTC1Url := "https://fullUrl1.com"
	var mockResTC1Err error = nil
	urlRepo.On(mockedGetUrlByID, ctx, mockArgTC1Id).Return(mockResTC1Url, mockResTC1Err).Once()
	argTC1ShortenUrl, err := encodeID(mockArgTC1Id)
	if err != nil {
		t.Fatal(err)
	}

	// Test case 2: URL is not found in the database
	mockArgTC2Id := int64(2)
	if err != nil {
		t.Fatal(err)
	}
	mockResTC2Url := ""
	mockResTC2Err := ErrUrlNotFound
	urlRepo.On(mockedGetUrlByID, ctx, mockArgTC2Id).Return(mockResTC2Url, mockResTC2Err).Once()
	argTC2ShortenUrl, err := encodeID(mockArgTC2Id)
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
			shortenUrl: argTC1ShortenUrl,
			want:       mockResTC1Url,
			wantErr:    false,
		},
		{
			name:       "URL is not found in the database",
			shortenUrl: argTC2ShortenUrl,
			want:       mockResTC2Url,
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
