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

func TestCreateTeam_Success(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	req := &dto.TeamRequest{
		Name:        tests.TestTeamName,
		Description: tests.TestTeamDescription,
		IsPublic:    tests.TestTeamIsPublic,
		UserId:      tests.TestUserID,
		TeamTopic:   tests.TestTeamTopic,
	}

	mockUserRepo.On("GetByID", tests.TestUserID).Return(&entity.User{ID: tests.TestUserID}, nil)
	mockRepo.On("Create", mock.AnythingOfType("*entity.Team")).Return(nil)
	mockRepo.On("GetTeamById", mock.Anything).Return(&tests.ValidTeam, nil)

	team, err := ts.CreateTeam(req)
	assert.NoError(t, err)
	assert.Equal(t, tests.TestTeamName, team.Name)
}

func TestAddUserToTeam_Success(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	mockUser := &entity.User{ID: tests.TestUserID}
	mockTeam := &entity.Team{Id: tests.TestTeamID, UsersIds: []string{}}

	mockUserRepo.On("GetByID", tests.TestUserID).Return(mockUser, nil)
	mockRepo.On("GetTeamById", tests.TestTeamID).Return(mockTeam, nil)
	mockUserRepo.On("Update", mockUser).Return(nil)
	mockRepo.On("Update", mockTeam).Return(nil)

	user, team, err := ts.AddUserToTeam(tests.TestUserID, tests.TestTeamID)
	assert.NoError(t, err)
	assert.Contains(t, team.UsersIds, tests.TestUserID)
	assert.Contains(t, *user.TeamsIds, tests.TestTeamID)
}

func TestAddUserToTeam_AlreadyMember(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	mockTeam := &entity.Team{Id: tests.TestTeamID, UsersIds: []string{tests.TestUserID}}
	mockUser := &entity.User{ID: tests.TestUserID}

	mockUserRepo.On("GetByID", tests.TestUserID).Return(mockUser, nil)
	mockRepo.On("GetTeamById", tests.TestTeamID).Return(mockTeam, nil)

	_, _, err := ts.AddUserToTeam(tests.TestUserID, tests.TestTeamID)
	assert.Error(t, err)
	assert.Equal(t, "user is already part of the team", err.Error())
}
