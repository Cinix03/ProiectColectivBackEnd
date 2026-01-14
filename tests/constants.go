package tests

const (
	// Test User Data
	TestFirstName    = "John"
	TestLastName     = "Doe"
	TestUsername     = "johndoe"
	TestEmail        = "john@example.com"
	TestPassword     = "password123"
	ExistingUsername = "existinguser"
	ExistingEmail    = "existing@example.com"

	// Test IDs
	TestUserID  = "123"
	TestTeamID  = "team123"
	TestTeamID2 = "team456"
	TestUserID1 = "user1"
	TestUserID2 = "user2"

	// Team Test Data
	TestTeamName        = "Awesome Team"
	TestTeamDescription = "This is a test team"
	TestTeamTopic       = "Programming"
	TestTeamIsPublic    = true

	// Team Request Test Data
	TestTeamRequestID = "teamRequest123"

	// Error Messages
	ErrUserNotFound    = "user not found"
	ErrUsernameExists  = "username already exists"
	ErrEmailExists     = "email already exists"
	ErrInvalidDuration = "Invalid timeSpentOnApp format"

	// Success Messages
	MsgStatisticsUpdated = "Statistics updated successfully"

	// HTTP Methods and Paths
	HTTPMethodPOST     = "POST"
	HTTPMethodPUT      = "PUT"
	PathUsersSignup    = "/users/signup"
	PathUserStatistics = "/users/123/statistics"
	PathEvents         = "/events"

	// Content Types
	ContentTypeJSON = "application/json"

	// Test Duration Values (milliseconds)
	TestDurationApp  = int64(9000000) // 2h30m in milliseconds
	TestDurationTeam = int64(4500000) // 1h15m in milliseconds

	// Gin Param Keys
	ParamKeyID = "id"

	// JSON Keys
	JSONKeyError   = "error"
	JSONKeyMessage = "message"

	// Test Event Data
	TestEventID          = "event123"
	TestEventName        = "Team Meeting"
	TestEventDescription = "Weekly sync meeting"
	TestEventDuration    = int64(3600000) // 1 hour in milliseconds

	// Event Error Messages
	ErrTeamNotFound       = "team not found"
	ErrEventNotFound      = "event not found"
	ErrInvalidEventStatus = "invalid event status"
	ErrStatusAlreadySet   = "status is already set to this value"
	ErrUserNotInEvent     = "user is not part of this event"
)
