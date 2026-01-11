package tests

import (
	"time"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
)

var (
	ValidSignUpRequest = dto.SignUpUserRequest{
		FirstName:        TestFirstName,
		LastName:         TestLastName,
		Username:         TestUsername,
		Email:            TestEmail,
		Password:         TestPassword,
		TopicsOfInterest: &[]model.TopicOfInterest{model.Programming},
	}

	ExistingUsernameRequest = dto.SignUpUserRequest{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Username:  ExistingUsername,
		Email:     TestEmail,
		Password:  TestPassword,
	}

	ExistingEmailRequest = dto.SignUpUserRequest{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Username:  TestUsername,
		Email:     ExistingEmail,
		Password:  TestPassword,
	}

	ValidSignUpResponse = dto.SignUpUserResponse{
		FirstName: TestFirstName,
		LastName:  TestLastName,
		Username:  TestUsername,
	}

	ExistingUser = entity.User{
		Username: ExistingUsername,
	}

	ValidUpdateStatisticsRequest = dto.UpdateStatisticsRequest{
		TimeSpentOnApp:  TestDurationApp,
		TeamId:          TestTeamID,
		TimeSpentOnTeam: TestDurationTeam,
	}

	ValidTimeSpentOnTeam = model.TimeSpentOnTeam{
		TeamId:   TestTeamID,
		Duration: int64(4500000), // 75 minutes in milliseconds
	}

	ValidFriendRequest = entity.FriendRequest{
		FromUserID: TestUserID1,
		ToUserID:   TestUserID2,
		Status:     entity.PENDING,
	}

	AcceptedFriendRequest = entity.FriendRequest{
		FromUserID: TestUserID1,
		ToUserID:   TestUserID2,
		Status:     entity.ACCEPTED,
	}

	DeniedFriendRequest = entity.FriendRequest{
		FromUserID: TestUserID1,
		ToUserID:   TestUserID2,
		Status:     entity.DENIED,
	}

	TestEventStartsAt = time.Now().Add(24 * time.Hour)

	ValidCreateEventRequest = dto.CreateEventRequest{
		InitiatorID: TestUserID,
		TeamID:      TestTeamID,
		Name:        TestEventName,
		Description: TestEventDescription,
		StartsAt:    TestEventStartsAt.String(),
		Duration:    TestEventDuration,
	}

	ValidEvent = entity.Event{
		ID:          TestEventID,
		InitiatorID: TestUserID,
		TeamID:      TestTeamID,
		Name:        TestEventName,
		Description: TestEventDescription,
		CreatedAt:   time.Now(),
		StartsAt:    TestEventStartsAt,
		Duration:    TestEventDuration,
		Statuses: map[string]entity.EventStatus{
			TestUserID1: entity.StatusPending,
			TestUserID2: entity.StatusPending,
		},
	}

	ValidUpdateEventRequest = dto.UpdateEventRequest{
		Name:        "Updated Event Name",
		Description: "Updated Description",
		StartsAt:    TestEventStartsAt.Add(1 * time.Hour).String(),
		Duration:    TestEventDuration * 2,
	}

	ValidUpdateEventStatusRequest = dto.UpdateEventStatusRequest{
		UserID: TestUserID1,
		Status: "accepted",
	}
)

func GetTestEventStartTime() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func GetTestEventStartTimeString() string {
	return GetTestEventStartTime().Format(time.RFC3339)
}

func GetValidCreateEventRequest() dto.CreateEventRequest {
	return dto.CreateEventRequest{
		InitiatorID: TestUserID,
		TeamID:      TestTeamID,
		Name:        TestEventName,
		Description: TestEventDescription,
		StartsAt:    GetTestEventStartTimeString(),
		Duration:    TestEventDuration,
	}
}

func GetValidEvent() entity.Event {
	return entity.Event{
		ID:          TestEventID,
		InitiatorID: TestUserID,
		TeamID:      TestTeamID,
		Name:        TestEventName,
		Description: TestEventDescription,
		CreatedAt:   time.Now(),
		StartsAt:    GetTestEventStartTime(),
		Duration:    TestEventDuration,
		Statuses: map[string]entity.EventStatus{
			TestUserID1: entity.StatusPending,
			TestUserID2: entity.StatusPending,
		},
	}
}

func GetValidUpdateEventRequest() dto.UpdateEventRequest {
	return dto.UpdateEventRequest{
		Name:        "Updated Event Name",
		Description: "Updated Description",
		StartsAt:    GetTestEventStartTime().Add(1 * time.Hour).Format(time.RFC3339),
		Duration:    TestEventDuration * 2,
	}
}
