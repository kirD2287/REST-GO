package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kirD2287/REST-GO/pkg/service"
	mock_service "github.com/kirD2287/REST-GO/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdenrity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string) 
	testTable := []struct {
		name           string
        headerName    string
		headerValue string
		token string
        mockBehavior  mockBehavior
		expectedStatusCode int
		expectedResponseBody string
	} {
		{
			name:           "OK",
            headerName:    "Authorization",
            headerValue: "Bearer token",
            token: "token",
            mockBehavior: func(s *mock_service.MockAuthorization, token string) {
                s.EXPECT().ParseToken("token").Return(1, nil)
            },
            expectedStatusCode: 200,
            expectedResponseBody: "1",
		},
		{
			name:           "No Header",
            headerName:    "",
			headerValue: "",
            mockBehavior: func(s *mock_service.MockAuthorization, token string) {},
            expectedStatusCode: 401,
            expectedResponseBody: `{"status":"empty auth header"}`,
		},

		{
			name:           "Invalid Bearer",
            headerName:    "Authorization",
            headerValue: "Invalid token",
            token: "",
            mockBehavior: func(s *mock_service.MockAuthorization, token string) {},
            expectedStatusCode: 401,
            expectedResponseBody: `{"status":"invalid auth header"}`,
		},

		{
			name:           "Invalid Token",
            headerName:    "Authorization",
            headerValue: "Bearer ",
            mockBehavior: func(s *mock_service.MockAuthorization, token string) {},
            expectedStatusCode: 401,
            expectedResponseBody: `{"status":"token is empty"}`,
		},
		{
			name:           "Service Failure",
            headerName:    "Authorization",
            headerValue: "Bearer token",
            token: "token",
            mockBehavior: func(s *mock_service.MockAuthorization, token string) {
                s.EXPECT().ParseToken("token").Return(1, errors.New("failed to parse token"))
            },
            expectedStatusCode: 401,
            expectedResponseBody: `{"status":"failed to parse token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
            c := gomock.NewController(t)
            defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
            testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
            handler := NewHandler(services)

			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(c *gin.Context) {
				id, exists := c.Get(userCtx)
				if !exists {
					c.JSON(401, gin.H{"status":"invalid auth header"})
				}
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
            req := httptest.NewRequest("GET", "/protected", nil)
			
            if testCase.headerName!= "" {
                req.Header.Set(testCase.headerName, testCase.headerValue)
            }
            req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)
			
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
			
	})

}

}