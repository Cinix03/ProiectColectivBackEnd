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

func TestNewTeamController_NewTeam(t *testing.T) {
	mockService := &tests.MockTeamService{}
	ctrl := controller.NewTeamControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	reqBody := dto.TeamRequest{
		Name:        "Test Team",
		Description: "A test team",
		IsPublic:    true,
		UserId:      "user1",
		TeamTopic:   "Topic1",
	}
	body, _ := json.Marshal(reqBody)
	c.Request, _ = http.NewRequest("POST", "/teams", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	fakeTeam := &entity.Team{
		Id:          "team123",
		Name:        reqBody.Name,
		Description: reqBody.Description,
		IsPublic:    reqBody.IsPublic,
		UsersIds:    []string{reqBody.UserId},
		TeamTopic:   reqBody.TeamTopic,
	}

	mockService.On("CreateTeam", &reqBody).Return(fakeTeam, nil)

	ctrl.NewTeam(c)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestGetTeam_Success(t *testing.T) {
	mockService := &tests.MockTeamService{}
	ctrl := controller.NewTeamControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: "team123"}}

	fakeTeam := &entity.Team{
		Id:          "team123",
		Name:        "Test Team",
		Description: "A test team",
		IsPublic:    true,
		UsersIds:    []string{"user1"},
		TeamTopic:   "Topic1",
	}

	mockService.On("GetTeamById", "team123").Return(fakeTeam, nil)

	ctrl.GetTeam(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestAddUserToTeam_Success(t *testing.T) {
	mockService := &tests.MockTeamService{}
	ctrl := controller.NewTeamControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	req := dto.UserToTeamRequest{
		UserID: "user2",
		TeamID: "team123",
	}
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("PUT", "/teams/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	fakeUser := &entity.User{
		ID:        req.UserID,
		FirstName: "New User",
	}
	fakeTeam := &entity.Team{
		Id:       req.TeamID,
		Name:     "Test Team",
		UsersIds: []string{"user1", req.UserID},
	}

	mockService.On("AddUserToTeam", req.UserID, req.TeamID).Return(fakeUser, fakeTeam, nil)

	ctrl.AddUserToTeam(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteUserFromTeam_Success(t *testing.T) {
	mockService := &tests.MockTeamService{}
	ctrl := controller.NewTeamControllerWithService(mockService)
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	req := dto.UserToTeamRequest{
		UserID: "user2",
		TeamID: "team123",
	}
	body, _ := json.Marshal(req)
	c.Request, _ = http.NewRequest("DELETE", "/teams/users", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	fakeUser := &entity.User{
		ID:        req.UserID,
		FirstName: "Removed User",
	}
	fakeTeam := &entity.Team{
		Id:       req.TeamID,
		Name:     "Test Team",
		UsersIds: []string{"user1"},
	}

	mockService.On("DeleteUserFromTeam", req.UserID, req.TeamID).Return(fakeUser, fakeTeam, nil)

	ctrl.DeleteUserFromTeam(c)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}
