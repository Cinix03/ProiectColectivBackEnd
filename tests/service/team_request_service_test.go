package service_test

import (
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTeamRequest_Success(t *testing.T) {
	mockTRRepo := &tests.MockTeamRequestRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	mockTeamRepo := &tests.MockTeamRepository{}
	service := service.NewTeamRequestServiceWithRepo(mockTRRepo, mockUserRepo, mockTeamRepo, nil)

	req := &dto.TeamRequestCreateDTO{
		UserID: tests.TestUserID1,
		TeamID: tests.TestTeamID,
	}

	mockUserRepo.On("GetByID", req.UserID).Return(&entity.User{ID: req.UserID}, nil)
	mockTeamRepo.On("GetTeamById", req.TeamID).Return(&tests.ValidTeam, nil)
	mockTRRepo.On("GetByUserId", req.UserID).Return([]*entity.TeamRequest{}, nil)
	mockTRRepo.On("Create", mock.AnythingOfType("*entity.TeamRequest")).Return(nil)

	teamReq, err := service.CreateTeamRequest(req)
	assert.NoError(t, err)
	assert.Equal(t, req.UserID, teamReq.UserID)
	assert.Equal(t, req.TeamID, teamReq.TeamID)
}

func TestCreateTeamRequest_UserAlreadyMember(t *testing.T) {
	mockTRRepo := &tests.MockTeamRequestRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	mockTeamRepo := &tests.MockTeamRepository{}
	service := service.NewTeamRequestServiceWithRepo(mockTRRepo, mockUserRepo, mockTeamRepo, nil)

	req := &dto.TeamRequestCreateDTO{
		UserID: tests.TestUserID,
		TeamID: tests.TestTeamID,
	}

	mockUserRepo.On("GetByID", req.UserID).Return(&entity.User{ID: req.UserID}, nil)
	mockTeamRepo.On("GetTeamById", req.TeamID).Return(&entity.Team{Id: req.TeamID, UsersIds: []string{req.UserID}}, nil)

	_, err := service.CreateTeamRequest(req)
	assert.Error(t, err)
	assert.Equal(t, "user is already a member of this team", err.Error())
}

func TestAcceptTeamRequest_Success(t *testing.T) {
	mockTRRepo := &tests.MockTeamRequestRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	mockTeamRepo := &tests.MockTeamRepository{}
	mockTeamService := &tests.MockTeamService{}
	service := service.NewTeamRequestServiceWithRepo(mockTRRepo, mockUserRepo, mockTeamRepo, mockTeamService)

	mockTRRepo.On("GetById", tests.TestTeamRequestID).Return(&tests.ValidTeamRequest, nil)
	mockTeamService.On("AddUserToTeam", tests.TestUserID1, tests.TestTeamID).Return(&entity.User{ID: tests.TestUserID1}, &tests.ValidTeam, nil)
	mockTRRepo.On("Delete", tests.TestTeamRequestID).Return(nil)

	user, team, err := service.AcceptTeamRequest(tests.TestTeamRequestID)
	assert.NoError(t, err)
	assert.Equal(t, tests.TestUserID1, user.ID)
	assert.Equal(t, tests.ValidTeam.Id, team.Id)
}

func TestRejectTeamRequest_Success(t *testing.T) {
	mockTRRepo := &tests.MockTeamRequestRepository{}
	service := service.NewTeamRequestServiceWithRepo(mockTRRepo, nil, nil, nil)

	mockTRRepo.On("GetById", tests.TestTeamRequestID).Return(&tests.ValidTeamRequest, nil)
	mockTRRepo.On("Delete", tests.TestTeamRequestID).Return(nil)

	err := service.RejectTeamRequest(tests.TestTeamRequestID)
	assert.NoError(t, err)
}

func TestGetAllTeamRequests(t *testing.T) {
	mockTRRepo := &tests.MockTeamRequestRepository{}
	service := service.NewTeamRequestServiceWithRepo(mockTRRepo, nil, nil, nil)

	mockTRRepo.On("GetAll").Return([]*entity.TeamRequest{&tests.ValidTeamRequest}, nil)

	reqs, err := service.GetAll()
	assert.NoError(t, err)
	assert.Len(t, reqs, 1)
}

func TestGetByUserId(t *testing.T) {
	mockTRRepo := &tests.MockTeamRequestRepository{}
	service := service.NewTeamRequestServiceWithRepo(mockTRRepo, nil, nil, nil)

	mockTRRepo.On("GetByUserId", tests.TestUserID1).Return([]*entity.TeamRequest{&tests.ValidTeamRequest}, nil)

	reqs, err := service.GetByUserId(tests.TestUserID1)
	assert.NoError(t, err)
	assert.Len(t, reqs, 1)
	assert.Equal(t, tests.TestUserID1, reqs[0].UserID)
}
