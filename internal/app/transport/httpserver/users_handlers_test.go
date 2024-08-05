// В файле, где объявлена структура Server

package httpserver

import (
	"bytes"
	"encoding/json"
	"github.com/vlad19930514/webApp/internal/app/transport/httpserver/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vlad19930514/webApp/internal/app/domain"
	"github.com/vlad19930514/webApp/internal/app/services"
	"github.com/vlad19930514/webApp/internal/app/services/mock_services"
)

func setupRouter(mockUserRepo *mocks.MockUserRepository) *gin.Engine {
	server := &httpserver.Server{
		Services: services.Services{Users: mockUserRepo},
	}
	router := gin.Default()
	router.POST("/users", server.CreateUser)
	router.GET("/users/:id", server.GetUser)
	router.PUT("/users", server.UpdateUser)
	return router
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	router := setupRouter(mockUserRepo)

	newUser := domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}
	mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(newUser, nil)

	reqBody, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var createdUser domain.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, newUser, createdUser)
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_services.NewMockUserRepository(ctrl)
	router := setupRouter(mockUserRepo)

	id := uuid.New()
	user := domain.User{
		ID:        id,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Age:       30,
	}
	mockUserRepo.EXPECT().GetUser(gomock.Any(), id).Return(user, nil)

	req, _ := http.NewRequest("GET", "/users/"+id.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var fetchedUser domain.User
	err := json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, user, fetchedUser)
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_services.NewMockUserRepository(ctrl)
	router := setupRouter(mockUserRepo)

	id := uuid.New()
	updatedUser := domain.User{
		ID:        id,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Age:       28,
	}
	mockUserRepo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(updatedUser, nil)

	reqBody, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var user domain.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser, user)
}
