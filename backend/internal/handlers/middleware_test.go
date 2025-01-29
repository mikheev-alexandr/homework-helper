package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/mikheev-alexandr/pet-project/backend/internal/service"
	mock_service "github.com/mikheev-alexandr/pet-project/backend/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_TeacherIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		cookieName          string
		cookieValue         string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			cookieName:  "Authorization",
			cookieValue: "Bearer mocked_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("mocked_token").Return(1, 0, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"user_id":1}`,
		},
		{
			name:                "Invalid Cookie Format",
			cookieName:          "Authorization",
			cookieValue:         "invalid_token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid auth cookie"}`,
		},
		{
			name:                "Token Empty",
			cookieName:          "Authorization",
			cookieValue:         "Bearer ",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"token is empty"}`,
		},
		{
			name:        "Not Enough Permissions",
			cookieName:  "Authorization",
			cookieValue: "Bearer valid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("valid_token").Return(1, 1, nil) // role != 0
			},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"not enough permissions for your role"}`,
		},
		{
			name:        "Invalid Token",
			cookieName:  "Authorization",
			cookieValue: "Bearer invalid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("invalid_token").Return(0, 0, fmt.Errorf("invalid token"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid token"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repos := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repos, test.cookieValue)

			services := &service.Service{Authorization: repos}
			handler := Handler{service: services}

			checkUserID := func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if exists {
					c.JSON(http.StatusOK, gin.H{"user_id": userID})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
				}
			}

			r := gin.New()
			r.GET("/teacher-identity", handler.teacherIdentity, checkUserID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/teacher-identity", nil)
			req.AddCookie(&http.Cookie{
				Name:  test.cookieName,
				Value: test.cookieValue,
			})

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedRequestBody, w.Body.String())
		})
	}
}

func TestMiddleware_StudentIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		cookieName          string
		cookieValue         string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			cookieName:  "Authorization",
			cookieValue: "Bearer mocked_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("mocked_token").Return(1, 1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"user_id":1}`,
		},
		{
			name:                "Invalid Cookie Format",
			cookieName:          "Authorization",
			cookieValue:         "invalid_token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid auth cookie"}`,
		},
		{
			name:                "Token Empty",
			cookieName:          "Authorization",
			cookieValue:         "Bearer ",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"token is empty"}`,
		},
		{
			name:        "Not Enough Permissions",
			cookieName:  "Authorization",
			cookieValue: "Bearer valid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("valid_token").Return(1, 0, nil)
			},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"not enough permissions for your role"}`,
		},
		{
			name:        "Invalid Token",
			cookieName:  "Authorization",
			cookieValue: "Bearer invalid_token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken("invalid_token").Return(0, 0, fmt.Errorf("invalid token"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"error":"invalid token"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repos := mock_service.NewMockAuthorization(c)
			test.mockBehavior(repos, "mocked_token")

			services := &service.Service{Authorization: repos}
			handler := Handler{service: services}

			checkUserID := func(c *gin.Context) {
				userID, exists := c.Get("user_id")
				if exists {
					c.JSON(http.StatusOK, gin.H{"user_id": userID})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found"})
				}
			}

			r := gin.New()
			r.GET("/student-identity", handler.studentIdentity, checkUserID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/student-identity", nil)
			req.AddCookie(&http.Cookie{
				Name:  test.cookieName,
				Value: test.cookieValue,
			})

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedRequestBody, w.Body.String())
		})
	}
}
