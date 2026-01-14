package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateTeamRequest(t *testing.T) {
	mockService := &tests.MockTeamRequestService{}
	ctrl := controller.NewTeamRequestControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	req := dto.TeamRequestCreateDTO{
		UserID: "user1",
		TeamID: "team123",
	}
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("POST", "/teamRequests", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	fakeRequest := &entity.TeamRequest{
		Id:     "req123",
		UserID: req.UserID,
		TeamID: req.TeamID,
	}

	mockService.On("CreateTeamRequest", &req).Return(fakeRequest, nil)

	ctrl.CreateTeamRequest(c)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestAcceptTeamRequest(t *testing.T) {
	mockService := &tests.MockTeamRequestService{}
	ctrl := controller.NewTeamRequestControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "req123"}}

	fakeUser := &entity.User{
		ID:        "user1",
		FirstName: "Test User",
	}
	fakeTeam := &entity.Team{
		Id:       "team123",
		Name:     "Test Team",
		UsersIds: []string{"user1"},
	}

	mockService.On("AcceptTeamRequest", "req123").Return(fakeUser, fakeTeam, nil)

	ctrl.AcceptTeamRequest(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestRejectTeamRequest(t *testing.T) {
	mockService := &tests.MockTeamRequestService{}
	ctrl := controller.NewTeamRequestControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "req123"}}

	mockService.On("RejectTeamRequest", "req123").Return(nil)

	ctrl.RejectTeamRequest(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllTeamRequests(t *testing.T) {
	mockService := &tests.MockTeamRequestService{}
	ctrl := controller.NewTeamRequestControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	fakeRequests := []*entity.TeamRequest{
		{Id: "req1", UserID: "user1", TeamID: "team123"},
		{Id: "req2", UserID: "user2", TeamID: "team123"},
	}

	mockService.On("GetAll").Return(fakeRequests, nil)

	ctrl.GetAllTeamRequests(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetTeamRequestsByUser(t *testing.T) {
	mockService := &tests.MockTeamRequestService{}
	ctrl := controller.NewTeamRequestControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "userId", Value: "user1"}}

	fakeRequests := []*entity.TeamRequest{
		{Id: "req1", UserID: "user1", TeamID: "team123"},
	}

	mockService.On("GetByUserId", "user1").Return(fakeRequests, nil)

	ctrl.GetTeamRequestsByUser(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}
