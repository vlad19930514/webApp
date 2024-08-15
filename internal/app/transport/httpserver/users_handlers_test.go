package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/vlad19930514/webApp/internal/app/domain"
	"github.com/vlad19930514/webApp/internal/app/transport/httpserver/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type UserTestSuite struct {
	suite.Suite
	mCtrl           *gomock.Controller
	mockUserService *mocks.MockIUserService
	router          *gin.Engine
}

// SetupTest is called before each test in the suite
func (s *UserTestSuite) SetupTest() {
	s.mCtrl = gomock.NewController(s.T())
	s.mockUserService = mocks.NewMockIUserService(s.mCtrl)

	server := NewHttpServer(s.mockUserService)
	router := gin.Default()
	router.POST("/user", server.createUser)
	router.GET("/user/:id", server.getUser)
	router.PUT("/user", server.updateUser)
	s.router = router
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestCreateUser() {
	tests := []struct {
		name           string
		input          createUserRequest
		mockReturnUser domain.User
		mockReturnErr  error
		expectedStatus int
		expectCall     bool
	}{
		{
			name: "successful creation",
			input: createUserRequest{
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Age:       28,
			},
			mockReturnUser: domain.User{
				Id:        uuid.New(),
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Age:       28,
				CreatedAt: time.Now(),
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectCall:     true,
		},
		{
			name: "validation error - invalid email",
			input: createUserRequest{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe",
				Age:       30,
			},
			mockReturnUser: domain.User{},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name: "server error",
			input: createUserRequest{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
				Age:       25,
			},
			mockReturnUser: domain.User{},
			mockReturnErr:  errors.New("internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectCall:     true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.expectCall {
				s.mockUserService.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(tt.mockReturnUser, tt.mockReturnErr)
			}

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.expectedStatus, w.Code, "Expected status code to be %v", tt.expectedStatus)

			if tt.expectedStatus == http.StatusOK {
				var createdUser domain.User
				err := json.Unmarshal(w.Body.Bytes(), &createdUser)
				assert.NoError(s.T(), err, "Expected no error when unmarshaling response body")
				assert.Equal(s.T(), tt.mockReturnUser.Id, createdUser.Id, "Expected created user ID to match")
				assert.Equal(s.T(), tt.mockReturnUser.FirstName, createdUser.FirstName, "Expected created user FirstName to match")
				assert.Equal(s.T(), tt.mockReturnUser.LastName, createdUser.LastName, "Expected created user LastName to match")
				assert.Equal(s.T(), tt.mockReturnUser.Email, createdUser.Email, "Expected created user Email to match")
				assert.Equal(s.T(), tt.mockReturnUser.Age, createdUser.Age, "Expected created user Age to match")
				assert.WithinDuration(s.T(), tt.mockReturnUser.CreatedAt, createdUser.CreatedAt, time.Second, "Expected created time to be within 1 second")
			}
		})
	}
}
func (s *UserTestSuite) TestGetUser() {
	tests := []struct {
		name           string
		userID         string
		mockReturnUser domain.User
		mockReturnErr  error
		expectedStatus int
		expectCall     bool
	}{
		{
			name:   "successful retrieval",
			userID: uuid.New().String(),
			mockReturnUser: domain.User{
				Id:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Age:       30,
				CreatedAt: time.Now(),
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusAccepted,
			expectCall:     true,
		},
		{
			name:           "validation error - invalid UUID",
			userID:         "invalid-uuid",
			mockReturnUser: domain.User{},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name:           "user not found",
			userID:         uuid.New().String(),
			mockReturnUser: domain.User{},
			mockReturnErr:  errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectCall:     true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.expectCall {
				s.mockUserService.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Return(tt.mockReturnUser, tt.mockReturnErr)
			}

			req, _ := http.NewRequest("GET", "/user/"+tt.userID, nil)

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.expectedStatus, w.Code, "Expected status code to be %v", tt.expectedStatus)

			if tt.expectedStatus == http.StatusAccepted {
				var retrievedUser domain.User
				err := json.Unmarshal(w.Body.Bytes(), &retrievedUser)
				assert.NoError(s.T(), err, "Expected no error when unmarshaling response body")
				assert.Equal(s.T(), tt.mockReturnUser.Id, retrievedUser.Id, "Expected retrieved user ID to match")
				assert.Equal(s.T(), tt.mockReturnUser.FirstName, retrievedUser.FirstName, "Expected retrieved user FirstName to match")
				assert.Equal(s.T(), tt.mockReturnUser.LastName, retrievedUser.LastName, "Expected retrieved user LastName to match")
				assert.Equal(s.T(), tt.mockReturnUser.Email, retrievedUser.Email, "Expected retrieved user Email to match")
				assert.Equal(s.T(), tt.mockReturnUser.Age, retrievedUser.Age, "Expected retrieved user Age to match")
				assert.WithinDuration(s.T(), tt.mockReturnUser.CreatedAt, retrievedUser.CreatedAt, time.Second, "Expected created time to be within 1 second")
			}
		})
	}
}
func (s *UserTestSuite) TestUpdateUser() {
	tests := []struct {
		name           string
		input          updateUserRequest
		mockReturnUser domain.User
		mockReturnErr  error
		expectedStatus int
		expectCall     bool
	}{
		{
			name: "successful update",
			input: updateUserRequest{
				ID:        uuid.New(),
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Age:       28,
			},
			mockReturnUser: domain.User{
				Id:        uuid.New(),
				FirstName: "Alice",
				LastName:  "Johnson",
				Email:     "alice.johnson@example.com",
				Age:       28,
				CreatedAt: time.Now(),
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectCall:     true,
		},
		{
			name: "validation error - invalid email",
			input: updateUserRequest{
				ID:        uuid.New(),
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe",
				Age:       30,
			},
			mockReturnUser: domain.User{},
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
			expectCall:     false,
		},
		{
			name: "server error",
			input: updateUserRequest{
				ID:        uuid.New(),
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
				Age:       25,
			},
			mockReturnUser: domain.User{},
			mockReturnErr:  errors.New("internal server error"),
			expectedStatus: http.StatusInternalServerError,
			expectCall:     true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.expectCall {
				s.mockUserService.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(tt.mockReturnUser, tt.mockReturnErr)
			}

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tt.expectedStatus, w.Code, "Expected status code to be %v", tt.expectedStatus)

			if tt.expectedStatus == http.StatusOK {
				var updatedUser domain.User
				err := json.Unmarshal(w.Body.Bytes(), &updatedUser)
				assert.NoError(s.T(), err, "Expected no error when unmarshaling response body")
				assert.Equal(s.T(), tt.mockReturnUser.Id, updatedUser.Id, "Expected updated user ID to match")
				assert.Equal(s.T(), tt.mockReturnUser.FirstName, updatedUser.FirstName, "Expected updated user FirstName to match")
				assert.Equal(s.T(), tt.mockReturnUser.LastName, updatedUser.LastName, "Expected updated user LastName to match")
				assert.Equal(s.T(), tt.mockReturnUser.Email, updatedUser.Email, "Expected updated user Email to match")
				assert.Equal(s.T(), tt.mockReturnUser.Age, updatedUser.Age, "Expected updated user Age to match")
			}
		})
	}
}
