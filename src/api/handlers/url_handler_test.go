package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yourusername/urlshortener/src/models"
)

// MockURLService is a mock implementation of URLService
type MockURLService struct {
	mock.Mock
}

// CreateShortURL implements the URLService interface
func (m *MockURLService) CreateShortURL(longURL string, expiresAt *time.Time) (*models.URL, error) {
	args := m.Called(longURL, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.URL), args.Error(1)
}

// GetLongURL implements the URLService interface
func (m *MockURLService) GetLongURL(shortID string) (string, error) {
	args := m.Called(shortID)
	return args.String(0), args.Error(1)
}

// Helper function to create a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

func TestShortenURL(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		input         map[string]interface{}
		mockResponse  *models.URL
		mockError     error
		expectedCode  int
		expectedError bool
	}{
		{
			name: "Valid URL",
			input: map[string]interface{}{
				"url":            "https://www.google.com",
				"expiration_days": 30,
			},
			mockResponse: &models.URL{
				ShortID:   "abc123",
				LongURL:   "https://www.google.com",
				ExpiresAt: timePtr(time.Now().AddDate(0, 0, 30)),
			},
			mockError:     nil,
			expectedCode:  http.StatusOK,
			expectedError: false,
		},
		{
			name: "Invalid URL",
			input: map[string]interface{}{
				"url":            "not-a-url",
				"expiration_days": 30,
			},
			mockResponse:  nil,
			mockError:     nil,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockURLService)
			if !tt.expectedError {
				mockService.On("CreateShortURL", mock.Anything, mock.Anything).Return(tt.mockResponse, tt.mockError)
			}

			// Create handler
			handler := NewURLHandler(mockService)

			// Create test request
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.ShortenURL(c)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)

			if !tt.expectedError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse.ShortID, response["short_url"])
				assert.Equal(t, tt.mockResponse.LongURL, response["long_url"])
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestRedirectToLongURL(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		shortID       string
		mockLongURL   string
		mockError     error
		expectedCode  int
		expectedError bool
	}{
		{
			name:          "Valid Short ID",
			shortID:       "abc123",
			mockLongURL:   "https://www.google.com",
			mockError:     nil,
			expectedCode:  http.StatusMovedPermanently,
			expectedError: false,
		},
		{
			name:          "Invalid Short ID",
			shortID:       "",
			mockLongURL:   "",
			mockError:     nil,
			expectedCode:  http.StatusBadRequest,
			expectedError: true,
		},
		{
			name:          "Non-existent Short ID",
			shortID:       "nonexistent",
			mockLongURL:   "",
			mockError:     assert.AnError,
			expectedCode:  http.StatusNotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockURLService)
			if tt.shortID != "" {
				mockService.On("GetLongURL", tt.shortID).Return(tt.mockLongURL, tt.mockError)
			}

			// Create handler
			handler := NewURLHandler(mockService)

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/"+tt.shortID, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "shortID", Value: tt.shortID}}

			// Call handler
			handler.RedirectToLongURL(c)

			// Assert response
			assert.Equal(t, tt.expectedCode, w.Code)

			if !tt.expectedError {
				assert.Equal(t, tt.mockLongURL, w.Header().Get("Location"))
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
} 