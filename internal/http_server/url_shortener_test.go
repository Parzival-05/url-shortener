package http_server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/Parzival-05/url-shortener/internal/http_server/io_server"
	"github.com/Parzival-05/url-shortener/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type UrlShortenerMock struct {
	mock.Mock
}

func (m *UrlShortenerMock) GetShortenUrl(ctx context.Context, fullUrl string) (string, error) {
	arg := m.Called(ctx, fullUrl)
	return arg.String(0), arg.Error(1)
}

func (m *UrlShortenerMock) SaveShortenUrl(ctx context.Context, fullUrl string) error {
	arg := m.Called(ctx, fullUrl)
	return arg.Error(0)
}

func (m *UrlShortenerMock) GetFullUrl(ctx context.Context, shortenUrl string) (string, error) {
	arg := m.Called(ctx, shortenUrl)
	return arg.String(0), arg.Error(1)
}

func (m *UrlShortenerMock) CreateUrl(ctx context.Context, fullUrl string) (string, error) {
	arg := m.Called(ctx, fullUrl)
	return arg.String(0), arg.Error(1)
}

func structToMapJSON(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func TestServer_CreateUrl(t *testing.T) {
	mockLog := zaptest.NewLogger(t)

	createPostRequest := func(req io_server.CreateUrlRequest) (*httptest.ResponseRecorder, *http.Request) {
		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/shorten", bytes.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		return w, r
	}

	ctx := context.Background()
	urlShortener := new(UrlShortenerMock)
	createUrl := "CreateUrl"
	// Test case 1: URL exists
	mockArgTC1Url := "https://fullUrl1.com"
	mockResTC1ShortUrl := "abc123"
	urlShortener.On(createUrl, ctx, mockArgTC1Url).Return(mockResTC1ShortUrl, nil)

	arg1 := io_server.CreateUrlRequest{URL: mockArgTC1Url}
	w1, r1 := createPostRequest(arg1)

	// Test case 2: URL not exists, save and get
	mockArgTC2Url := "https://fullUrl2.com"
	mockResTC2ShortUrl := "def456"
	urlShortener.On(createUrl, ctx, mockArgTC2Url).Return(mockResTC2ShortUrl, nil).Once()

	arg2 := io_server.CreateUrlRequest{URL: mockArgTC2Url}
	w2, r2 := createPostRequest(arg2)

	// Test case 3: Error on save
	mockArgTC3Url := "https://fullUrl3.com"
	urlShortener.On(createUrl, ctx, mockArgTC3Url).Return("", errors.New("Failed to save shorten url")).Once()

	arg3 := io_server.CreateUrlRequest{URL: mockArgTC3Url}
	w3, r3 := createPostRequest(arg3)

	server := Server{
		log:          mockLog,
		urlShortener: urlShortener,
	}

	tests := []struct {
		name     string
		w        *httptest.ResponseRecorder
		r        *http.Request
		expected struct {
			code int
			resp io_server.CreateUrlResponse
			err  string
		}
	}{
		{
			name: "URL already exists",
			w:    w1,
			r:    r1,
			expected: struct {
				code int
				resp io_server.CreateUrlResponse
				err  string
			}{
				code: http.StatusOK,
				resp: io_server.CreateUrlResponse{ShortenURL: mockResTC1ShortUrl},
			},
		},
		{
			name: "URL does not exist - successfully created",
			w:    w2,
			r:    r2,
			expected: struct {
				code int
				resp io_server.CreateUrlResponse
				err  string
			}{
				code: http.StatusOK,
				resp: io_server.CreateUrlResponse{ShortenURL: mockResTC2ShortUrl},
			},
		},
		{
			name: "Error saving URL",
			w:    w3,
			r:    r3,
			expected: struct {
				code int
				resp io_server.CreateUrlResponse
				err  string
			}{
				code: http.StatusInternalServerError,
				err:  "Failed to save shorten url",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.CreateUrl(tt.w, tt.r)

			assert.Equal(t, tt.expected.code, tt.w.Code)

			var response map[string]any
			err := json.Unmarshal(tt.w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expected.err != "" {
				assert.Contains(t, response["error"], tt.expected.err)
			} else if tt.expected.code == http.StatusOK {
				data, exists := response["data"]
				assert.True(t, exists)

				dataBytes, err := json.Marshal(data)
				assert.NoError(t, err)
				var actualResp io_server.CreateUrlResponse
				err = json.Unmarshal(dataBytes, &actualResp)
				assert.NoError(t, err)

				assert.Equal(t, tt.expected.resp.ShortenURL, actualResp.ShortenURL)
			}
		})
	}
}

func TestServer_GetUrl(t *testing.T) {
	mockLog := zaptest.NewLogger(t)
	getRW := func(req io_server.GetUrlRequest) (w *httptest.ResponseRecorder, r *http.Request) {
		params := url.Values{}
		params.Add("shorten_url", req.ShortenURL)
		reqURL := url.URL{
			Path:     "/shorten",
			RawQuery: params.Encode(),
		}
		w = httptest.NewRecorder()
		r = httptest.NewRequest("GET", reqURL.String(), nil)
		return w, r
	}

	ctx := context.Background()
	mockedGetFullUrl := "GetFullUrl"
	urlShortener := new(UrlShortenerMock)

	// Test case 1: URL is found in the database
	mockArgTC1Url := "shortenUrl"
	mockResTC1Url := "https://fullUrl1.com"
	var mockResTC1Err error = nil
	urlShortener.On(mockedGetFullUrl, ctx, mockArgTC1Url).Return(mockResTC1Url, mockResTC1Err).Once()
	arg := io_server.GetUrlRequest{
		ShortenURL: mockArgTC1Url,
	}
	w1, r1 := getRW(arg)
	res1 := io_server.GetUrlResponse{
		URL: mockResTC1Url,
	}
	var err1 error = nil
	code1 := http.StatusOK

	// Test case 2: URL is not found in the database
	mockArgTC2Url := "notShortenUrl"
	mockResTC2Url := ""
	mockResTC2Err := service.ErrUrlNotFound
	urlShortener.On(mockedGetFullUrl, ctx, mockArgTC2Url).Return(mockResTC2Url, mockResTC2Err).Once()
	arg = io_server.GetUrlRequest{
		ShortenURL: mockArgTC2Url,
	}
	w2, r2 := getRW(arg)
	res2 := io_server.GetUrlResponse{
		URL: mockResTC2Url,
	}
	err2 := service.ErrUrlNotFound
	code2 := http.StatusBadRequest

	server := Server{
		port:         0,
		log:          mockLog,
		db:           nil,
		urlShortener: urlShortener,
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		w *httptest.ResponseRecorder
		r *http.Request

		expected struct {
			res  io_server.GetUrlResponse
			code int
			err  error
		}
	}{
		{
			name: "Test case 1: URL is found in the database",
			w:    w1,
			r:    r1,
			expected: struct {
				res  io_server.GetUrlResponse
				code int
				err  error
			}{
				res:  res1,
				code: code1,
				err:  err1,
			},
		},
		{
			name: "Test case 2: URL is not found in the database",
			w:    w2,
			r:    r2,
			expected: struct {
				res  io_server.GetUrlResponse
				code int
				err  error
			}{
				res:  res2,
				code: code2,
				err:  err2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server.GetUrl(tt.w, tt.r)

			contentBytes := tt.w.Body.Bytes()
			var content map[string]any
			err := json.Unmarshal(contentBytes, &content)
			assert.NoError(t, err)

			code := tt.w.Result().StatusCode
			assert.Equal(t, tt.expected.code, code)

			if tt.expected.err != nil {
				errRes := content["error"].(string)
				assert.Equal(t, tt.expected.err.Error(), errRes)
				return
			}
			data, _ := content["data"].(map[string]any)

			res, err := structToMapJSON(tt.expected.res)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(res, data) {
				t.Errorf("Expected result %v, got %v", res, data)
			}
		})
	}
}
