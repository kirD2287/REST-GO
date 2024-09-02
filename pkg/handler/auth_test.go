package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	todo "github.com/kirD2287/REST-GO"
	"github.com/kirD2287/REST-GO/pkg/service"
	mock_service "github.com/kirD2287/REST-GO/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)
	testTable := []struct {
		name           string
		inputBody string
		inputUser todo.User
        mockBehavior   mockBehavior
		expectedStatusCode int
		expectedRequestBody string
        
	} {
		{
			name: "OK",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
                Username: "test",
                Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode: 200,
            expectedRequestBody: `{"id":1}`,
		},
		{
			name: "Empty Fields",
			inputBody: `{"name": "", "username": "", "password": ""}`,
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(gomock.Any()).Times(0)
			},
			expectedStatusCode: 400,
            expectedRequestBody: `{"status":"Invalid input body"}`,
		},

		{
			name: "Server Failure",
			inputBody: `{"name": "Test", "username": "test", "password": "qwerty"}`,
			inputUser: todo.User{
				Name:     "Test",
                Username: "test",
                Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("server failure"))
			},

			expectedStatusCode: 500,
            expectedRequestBody: `{"status":"server failure"}`,
		
	},
		
}
    for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
            defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
		bytes.NewBufferString(testCase.inputBody))

		    r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
}


}