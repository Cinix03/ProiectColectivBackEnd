package service_test

import (
	"fmt"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventService_CreateEvent_Success(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests. MockUserRepository)
	
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	teamMembers := []string{tests.TestUserID1, tests. TestUserID2}
	team := &entity.Team{
		Id:       tests.TestTeamID,
		UsersIds:  teamMembers,
	}
	user := &entity.User{ID: tests.TestUserID}

	request := tests.GetValidCreateEventRequest() // Use function instead of variable

	mockUserRepo.On("GetByID", tests.TestUserID).Return(user, nil)
	mockTeamRepo.On("GetTeamById", tests.TestTeamID).Return(team, nil)
	mockEventRepo.On("Create", mock. MatchedBy(func(e *entity.Event) bool {
		return e.Name == request.Name &&
			e.TeamID == request.TeamID &&
			e.InitiatorID == request. InitiatorID &&
			len(e.Statuses) == len(teamMembers)
	})).Return(nil)

	resp, err := es.CreateEvent(&request) // Pass pointer

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, request. Name, resp.Name)
	assert.Equal(t, request. TeamID, resp.TeamID)
	assert.NotEmpty(t, resp.ID)

	mockEventRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestEventService_CreateEvent_TeamNotFound(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	request := tests.GetValidCreateEventRequest() // Use function
	user := &entity.User{ID:  tests.TestUserID}

	mockUserRepo.On("GetByID", tests.TestUserID).Return(user, nil)
	mockTeamRepo.On("GetTeamById", tests.TestTeamID).Return(nil, fmt.Errorf(tests.ErrTeamNotFound))

	resp, err := es.CreateEvent(&request)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), tests.ErrTeamNotFound)

	mockEventRepo.AssertNotCalled(t, "Create", mock. Anything)
}

func TestEventService_CreateEvent_UserNotFound(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests. MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	request := tests.GetValidCreateEventRequest() // Use function

	mockUserRepo.On("GetByID", tests.TestUserID).Return(nil, fmt.Errorf("user not found"))

	resp, err := es.CreateEvent(&request)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "user not found")

	mockEventRepo.AssertNotCalled(t, "Create", mock.Anything)
	mockTeamRepo.AssertNotCalled(t, "GetTeamById", mock.Anything)
}

func TestEventService_GetEventById_Success(t *testing. T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	event := tests.GetValidEvent() // Use function

	mockEventRepo. On("GetByID", tests.TestEventID).Return(&event, nil)

	resp, err := es.GetEventById(tests.TestEventID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, event.ID, resp.ID)
	assert.Equal(t, event.Name, resp.Name)

	mockEventRepo.AssertExpectations(t)
}

func TestEventService_GetEventById_NotFound(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	mockEventRepo.On("GetByID", "invalid-id").Return(nil, fmt.Errorf(tests.ErrEventNotFound))

	resp, err := es. GetEventById("invalid-id")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), tests.ErrEventNotFound)

	mockEventRepo.AssertExpectations(t)
}

func TestEventService_GetEventsByTeamId_Success(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	event1 := tests.GetValidEvent()
	event2 := tests. GetValidEvent()
	event2.ID = "event456"
	events := []*entity.Event{&event1, &event2}

	mockEventRepo.On("GetByTeamID", tests.TestTeamID).Return(events, nil)

	resp, err := es.GetEventsByTeamId(tests.TestTeamID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp, 2)

	mockEventRepo.AssertExpectations(t)
}

func TestEventService_UpdateEventDetails_Success(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	existingEvent := tests.GetValidEvent()
	request := tests.GetValidUpdateEventRequest() // Use function

	mockEventRepo.On("GetByID", tests.TestEventID).Return(&existingEvent, nil)
	mockEventRepo.On("Update", tests.TestEventID, mock. MatchedBy(func(updates map[string]interface{}) bool {
		return updates["name"] == request.Name &&
			updates["description"] == request.Description
	})).Return(nil)

	resp, err := es.UpdateEventDetails(tests.TestEventID, &request)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, request.Name, resp.Name)
	assert.Equal(t, request.Description, resp.Description)

	mockEventRepo.AssertExpectations(t)
}

func TestEventService_UpdateEventDetails_EventNotFound(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	request := tests.GetValidUpdateEventRequest() // Use function

	mockEventRepo.On("GetByID", "invalid-id").Return(nil, fmt.Errorf(tests. ErrEventNotFound))

	resp, err := es.UpdateEventDetails("invalid-id", &request)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockEventRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestEventService_UpdateUserStatus_Success(t *testing. T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	existingEvent := tests.GetValidEvent()
	request := &tests.ValidUpdateEventStatusRequest
	user := &entity.User{ID:  request.UserID}

	mockEventRepo.On("GetByID", tests.TestEventID).Return(&existingEvent, nil)
	mockUserRepo.On("GetByID", request.UserID).Return(user, nil)
	mockEventRepo.On("Update", tests. TestEventID, mock.MatchedBy(func(updates map[string]interface{}) bool {
		if statuses, ok := updates["statuses"].(map[string]entity.EventStatus); ok {
			return statuses[request.UserID] == entity.EventStatus(request.Status)
		}
		return false
	})).Return(nil)

	resp, err := es.UpdateUserStatus(tests.TestEventID, request)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockEventRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestEventService_UpdateUserStatus_SameStatusConflict(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	existingEvent := tests.GetValidEvent()
	existingEvent.Statuses[tests.TestUserID1] = entity.StatusAccepted

	request := &dto.UpdateEventStatusRequest{
		UserID: tests.TestUserID1,
		Status: string(entity.StatusAccepted), // Same as current
	}
	user := &entity.User{ID: request.UserID}

	mockEventRepo.On("GetByID", tests.TestEventID).Return(&existingEvent, nil)
	mockUserRepo.On("GetByID", request.UserID).Return(user, nil)

	resp, err := es.UpdateUserStatus(tests.TestEventID, request)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), service.SameStatus)

	mockEventRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestEventService_UpdateUserStatus_InvalidStatus(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	existingEvent := tests.GetValidEvent()

	request := &dto.UpdateEventStatusRequest{
		UserID:  tests.TestUserID1,
		Status: "invalid-status",
	}
	user := &entity.User{ID: request. UserID}

	mockEventRepo.On("GetByID", tests.TestEventID).Return(&existingEvent, nil)
	mockUserRepo.On("GetByID", request.UserID).Return(user, nil)

	resp, err := es.UpdateUserStatus(tests.TestEventID, request)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), service.InvalidStatus)

	mockEventRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestEventService_UpdateUserStatus_UserNotFound(t *testing.T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	existingEvent := tests.GetValidEvent()
	request := &tests.ValidUpdateEventStatusRequest

	mockEventRepo.On("GetByID", tests.TestEventID).Return(&existingEvent, nil)
	mockUserRepo.On("GetByID", request.UserID).Return(nil, fmt.Errorf("user not found"))

	resp, err := es.UpdateUserStatus(tests.TestEventID, request)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockEventRepo.AssertNotCalled(t, "Update", mock.Anything, mock. Anything)
}

func TestEventService_DeleteEvent_Success(t *testing.T) {
	mockEventRepo := new(tests. MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	event := tests.GetValidEvent()

	mockEventRepo.On("GetByID", tests.TestEventID).Return(&event, nil)
	mockEventRepo.On("Delete", tests.TestEventID).Return(nil)

	err := es.DeleteEvent(tests. TestEventID)

	assert.NoError(t, err)

	mockEventRepo.AssertExpectations(t)
}

func TestEventService_DeleteEvent_NotFound(t *testing. T) {
	mockEventRepo := new(tests.MockEventRepository)
	mockTeamRepo := new(tests.MockTeamRepository)
	mockUserRepo := new(tests.MockUserRepository)
	es := service.NewEventServiceWithRepo(mockEventRepo, mockTeamRepo, mockUserRepo)

	mockEventRepo.On("GetByID", "invalid-id").Return(nil, fmt.Errorf(tests. ErrEventNotFound))

	err := es.DeleteEvent("invalid-id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), tests.ErrEventNotFound)

	mockEventRepo.AssertNotCalled(t, "Delete", mock.Anything)
}
