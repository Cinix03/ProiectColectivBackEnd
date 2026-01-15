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

	testUser := &entity.User{
		ID:       tests.TestUserID,
		TeamsIds: &[]string{},
	}

	emptyTeam := &entity.Team{
		Id:          "test-team-id",
		Name:        tests.TestTeamName,
		Description: tests.TestTeamDescription,
		IsPublic:    tests.TestTeamIsPublic,
		UsersIds:    []string{},
		TeamTopic:   tests.TestTeamTopic,
	}

	teamWithUser := &entity.Team{
		Id:          "test-team-id",
		Name:        tests.TestTeamName,
		Description: tests.TestTeamDescription,
		IsPublic:    tests.TestTeamIsPublic,
		UsersIds:    []string{tests.TestUserID},
		TeamTopic:   tests.TestTeamTopic,
	}

	mockUserRepo.On("GetByID", tests.TestUserID).Return(testUser, nil)
	mockRepo.On("Create", mock.AnythingOfType("*entity.Team")).Return(nil)

	mockRepo.On("GetTeamById", mock.Anything).Return(emptyTeam, nil).Once()
	mockUserRepo.On("Update", mock.AnythingOfType("*entity.User")).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("*entity.Team")).Return(nil)

	mockRepo.On("GetTeamById", mock.Anything).Return(teamWithUser, nil)

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

func TestGetUsersByTeam_Success(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	user1 := &entity.User{
		ID:        "user1",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
	}
	user2 := &entity.User{
		ID:        "user2",
		FirstName: "Jane",
		LastName:  "Smith",
		Username:  "janesmith",
		Email:     "jane@example.com",
	}

	testTeam := &entity.Team{
		Id:       tests.TestTeamID,
		Name:     tests.TestTeamName,
		UsersIds: []string{"user1", "user2"},
	}

	mockRepo.On("GetTeamById", tests.TestTeamID).Return(testTeam, nil)
	mockUserRepo.On("GetByID", "user1").Return(user1, nil)
	mockUserRepo.On("GetByID", "user2").Return(user2, nil)

	users, err := ts.GetUsersByTeam(tests.TestTeamID)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "johndoe", users[0].Username)
	assert.Equal(t, "janesmith", users[1].Username)
}

func TestGetUsersByTeam_EmptyTeamID(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	users, err := ts.GetUsersByTeam("")
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, "team ID is required", err.Error())

	users, err = ts.GetUsersByTeam("   ")
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, "team ID is required", err.Error())
}

func TestGetUsersByTeam_TeamNotFound(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	mockRepo.On("GetTeamById", "non-existent-team").Return(nil, assert.AnError)

	users, err := ts.GetUsersByTeam("non-existent-team")
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Equal(t, "team not found", err.Error())
}

func TestGetUsersByTeam_EmptyTeam(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	// Create test team with no users
	testTeam := &entity.Team{
		Id:       tests.TestTeamID,
		Name:     tests.TestTeamName,
		UsersIds: []string{},
	}

	mockRepo.On("GetTeamById", tests.TestTeamID).Return(testTeam, nil)

	users, err := ts.GetUsersByTeam(tests.TestTeamID)
	assert.NoError(t, err)
	assert.Len(t, users, 0)
}

func TestGetUsersByTeam_WithInvalidUser(t *testing.T) {
	mockRepo := &tests.MockTeamRepository{}
	mockUserRepo := &tests.MockUserRepository{}
	ts := service.NewTeamServiceWithRepo(mockUserRepo, mockRepo)

	validUser := &entity.User{
		ID:        "user1",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Email:     "john@example.com",
	}

	testTeam := &entity.Team{
		Id:       tests.TestTeamID,
		Name:     tests.TestTeamName,
		UsersIds: []string{"user1", "invalid-user"},
	}

	mockRepo.On("GetTeamById", tests.TestTeamID).Return(testTeam, nil)
	mockUserRepo.On("GetByID", "user1").Return(validUser, nil)
	mockUserRepo.On("GetByID", "invalid-user").Return(nil, assert.AnError)

	users, err := ts.GetUsersByTeam(tests.TestTeamID)
	assert.NoError(t, err)  // Should not error, but skip invalid user
	assert.Len(t, users, 1) // Only valid user should be returned
	assert.Equal(t, "johndoe", users[0].Username)
}
