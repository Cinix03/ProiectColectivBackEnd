package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SerbanEduard/ProiectColectivBackEnd/controller"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/service"
	"github.com/SerbanEduard/ProiectColectivBackEnd/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEventController_NewEvent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := tests.GetValidCreateEventRequest()
	expectedResp := &dto.EventDTO{
		ID:          tests.TestEventID,
		InitiatorID: request.InitiatorID,
		TeamID:      request.TeamID,
		Name:        request.Name,
		Description: request.Description,
		StartsAt:    request.StartsAt,
		Duration:    request.Duration,
	}

	mockService.On("CreateEvent", &request).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(tests.HTTPMethodPOST, tests.PathEvents, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.NewEvent(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp dto.EventDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp.ID, resp.ID)
	assert.Equal(t, expectedResp.Name, resp.Name)

	mockService.AssertExpectations(t)
}

func TestEventController_NewEvent_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(tests.HTTPMethodPOST, tests.PathEvents, bytes.NewBuffer([]byte(`{invalid json}`)))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.NewEvent(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], "invalid")
}

func TestEventController_NewEvent_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := tests.GetValidCreateEventRequest()

	mockService.On("CreateEvent", &request).Return(nil, fmt.Errorf("service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(tests.HTTPMethodPOST, tests.PathEvents, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.NewEvent(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], "service error")

	mockService.AssertExpectations(t)
}

func TestEventController_GetEvent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	expectedEvent := &dto.EventDTO{
		ID:          tests.TestEventID,
		Name:        tests.TestEventName,
		Description: tests.TestEventDescription,
		TeamID:      tests.TestTeamID,
	}

	mockService.On("GetEventById", tests.TestEventID).Return(expectedEvent, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	ec.GetEvent(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.EventDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedEvent.ID, resp.ID)
	assert.Equal(t, expectedEvent.Name, resp.Name)

	mockService.AssertExpectations(t)
}

func TestEventController_GetEvent_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	mockService.On("GetEventById", "invalid-id").Return(nil, fmt.Errorf(tests.ErrEventNotFound))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "invalid-id"}}

	ec.GetEvent(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], tests.ErrEventNotFound)

	mockService.AssertExpectations(t)
}

func TestEventController_GetEvents_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockEventService := new(tests.MockEventService)
	mockTeamService := new(tests.MockTeamService)
	ec := controller.NewEventController()
	ec.SetEventService(mockEventService)
	ec.SetTeamService(mockTeamService)

	events := []*dto.EventDTO{
		{ID: tests.TestEventID, Name: tests.TestEventName, TeamID: tests.TestTeamID},
		{ID: "event456", Name: "Another Event", TeamID: tests.TestTeamID},
	}

	team := &entity.Team{Id: tests.TestTeamID}
	mockTeamService.On("GetTeamById", tests.TestTeamID).Return(team, nil)
	mockEventService.On("GetEventsByTeamId", tests.TestTeamID).Return(events, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, tests.PathEvents+"?teamId="+tests.TestTeamID, nil)

	ec.GetEvents(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []*dto.EventDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, events[0].ID, resp[0].ID)
	assert.Equal(t, events[1].ID, resp[1].ID)

	mockEventService.AssertExpectations(t)
	mockTeamService.AssertExpectations(t)
}

func TestEventController_GetEvents_MissingTeamID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, tests.PathEvents, nil)

	ec.GetEvents(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response[tests.JSONKeyError])
}

func TestEventController_GetEvents_TeamNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockEventService := new(tests.MockEventService)
	mockTeamService := new(tests.MockTeamService)
	ec := controller.NewEventController()
	ec.SetEventService(mockEventService)
	ec.SetTeamService(mockTeamService)

	mockTeamService.On("GetTeamById", "invalid-team").Return(nil, fmt.Errorf(tests.ErrTeamNotFound))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, tests.PathEvents+"?teamId=invalid-team", nil)

	ec.GetEvents(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], tests.ErrTeamNotFound)

	mockTeamService.AssertExpectations(t)
	mockEventService.AssertNotCalled(t, "GetEventsByTeamId")
}

func TestEventController_GetEvents_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockEventService := new(tests.MockEventService)
	mockTeamService := new(tests.MockTeamService)
	ec := controller.NewEventController()
	ec.SetEventService(mockEventService)
	ec.SetTeamService(mockTeamService)

	team := &entity.Team{Id: tests.TestTeamID}
	mockTeamService.On("GetTeamById", tests.TestTeamID).Return(team, nil)
	mockEventService.On("GetEventsByTeamId", tests.TestTeamID).Return(nil, fmt.Errorf("service error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, tests.PathEvents+"?teamId="+tests.TestTeamID, nil)

	ec.GetEvents(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockEventService.AssertExpectations(t)
	mockTeamService.AssertExpectations(t)
}

func TestEventController_UpdateEventDetails_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := tests.GetValidUpdateEventRequest()
	expectedResp := &dto.EventDTO{
		ID:          tests.TestEventID,
		Name:        request.Name,
		Description: request.Description,
		StartsAt:    request.StartsAt,
		Duration:    request.Duration,
	}

	mockService.On("UpdateEventDetails", tests.TestEventID, &request).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateEventDetails(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.EventDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp.Name, resp.Name)
	assert.Equal(t, expectedResp.Description, resp.Description)

	mockService.AssertExpectations(t)
}

func TestEventController_UpdateEventDetails_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID, bytes.NewBuffer([]byte(`{invalid json}`)))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateEventDetails(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventController_UpdateEventDetails_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := tests.GetValidUpdateEventRequest()

	mockService.On("UpdateEventDetails", tests.TestEventID, &request).Return(nil, fmt.Errorf(tests.ErrEventNotFound))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID, bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateEventDetails(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], tests.ErrEventNotFound)

	mockService.AssertExpectations(t)
}

func TestEventController_UpdateUserStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := &tests.ValidUpdateEventStatusRequest
	expectedResp := &dto.EventDTO{
		ID:            tests.TestEventID,
		AcceptedCount: 1,
		PendingCount:  1,
		DeclinedCount: 0,
	}

	mockService.On("UpdateUserStatus", tests.TestEventID, request).Return(expectedResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID+"/status", bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateUserStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp dto.EventDTO
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, expectedResp.ID, resp.ID)

	mockService.AssertExpectations(t)
}

func TestEventController_UpdateUserStatus_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID+"/status", bytes.NewBuffer([]byte(`{invalid json}`)))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateUserStatus(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestEventController_UpdateUserStatus_SameStatus_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := &tests.ValidUpdateEventStatusRequest

	mockService.On("UpdateUserStatus", tests.TestEventID, request).Return(nil, fmt.Errorf(service.SameStatus))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID+"/status", bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateUserStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], "already set")

	mockService.AssertExpectations(t)
}

func TestEventController_UpdateUserStatus_InvalidStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	request := &dto.UpdateEventStatusRequest{
		UserID: tests.TestUserID1,
		Status: "invalid-status",
	}

	mockService.On("UpdateUserStatus", tests.TestEventID, request).Return(nil, fmt.Errorf(service.InvalidStatus))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	jsonData, _ := json.Marshal(request)
	c.Request, _ = http.NewRequest(http.MethodPatch, "/events/"+tests.TestEventID+"/status", bytes.NewBuffer(jsonData))
	c.Request.Header.Set(tests.ContentTypeJSON, tests.ContentTypeJSON)

	ec.UpdateUserStatus(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], service.InvalidStatus)

	mockService.AssertExpectations(t)
}

func TestEventController_DeleteEvent_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	mockService.On("DeleteEvent", tests.TestEventID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: tests.TestEventID}}

	ec.DeleteEvent(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyMessage], "deleted")

	mockService.AssertExpectations(t)
}

func TestEventController_DeleteEvent_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(tests.MockEventService)
	ec := controller.NewEventController()
	ec.SetEventService(mockService)

	mockService.On("DeleteEvent", "invalid-id").Return(fmt.Errorf(tests.ErrEventNotFound))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "invalid-id"}}

	ec.DeleteEvent(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response[tests.JSONKeyError], tests.ErrEventNotFound)

	mockService.AssertExpectations(t)
}
