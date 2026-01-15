# Repository Summary: ProiectColectivBackEnd (StudyWithMe API)

## Overview

**ProiectColectivBackEnd** is a backend API server for "StudyWithMe" (also known as "StudyFlow"), a collaborative learning platform designed to facilitate team-based studying and knowledge sharing. The application is built with **Go** using the **Gin web framework** and **Firebase** as the database backend.

The API serves as the backbone for a study collaboration platform that enables users to create teams, share files, take quizzes, communicate in real-time, and participate in voice/video rooms for collaborative learning sessions.

## Technology Stack

### Core Technologies
- **Language**: Go 1.25.0
- **Web Framework**: Gin (v1.11.0)
- **Database**: Firebase Realtime Database & Firestore
- **Authentication**: JWT (JSON Web Tokens) using golang-jwt/v5
- **Real-time Communication**: WebSockets (gorilla/websocket v1.5.3)
- **API Documentation**: Swagger/OpenAPI (swaggo)

### Key Dependencies
- **firebase.google.com/go/v4**: Firebase Admin SDK for Go
- **github.com/gin-gonic/gin**: HTTP web framework
- **github.com/golang-jwt/jwt/v5**: JWT authentication
- **github.com/gorilla/websocket**: WebSocket implementation
- **github.com/swaggo/gin-swagger**: Swagger documentation
- **golang.org/x/crypto**: Password hashing and encryption
- **google.golang.org/api**: Google Cloud APIs

## Architecture

The application follows a **layered architecture** pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────┐
│                   Routes Layer                      │
│          (HTTP routing and middleware)              │
├─────────────────────────────────────────────────────┤
│                 Controller Layer                    │
│        (Request handling and validation)            │
├─────────────────────────────────────────────────────┤
│                  Service Layer                      │
│             (Business logic)                        │
├─────────────────────────────────────────────────────┤
│               Persistence Layer                     │
│         (Database interactions)                     │
├─────────────────────────────────────────────────────┤
│                  Model Layer                        │
│        (Entities, DTOs, Types)                      │
└─────────────────────────────────────────────────────┘
```

### Directory Structure

- **`/controller`**: HTTP request handlers and business logic coordination
- **`/service`**: Core business logic and orchestration
- **`/persistence`**: Repository pattern implementation for database operations
- **`/model`**: Data models, entities, and DTOs
- **`/routes`**: API route definitions and middleware setup
- **`/config`**: Configuration management (Firebase, JWT)
- **`/utils`**: Utility functions (JWT helpers)
- **`/validator`**: Input validation logic
- **`/hub`**: WebSocket hub for real-time messaging
- **`/tests`**: Test files and test utilities
- **`/docs`**: Auto-generated Swagger documentation

## Core Features

### 1. User Management
- User registration and authentication
- User profile management (CRUD operations)
- Password hashing with bcrypt
- JWT-based authentication

### 2. Team Collaboration
- Create and manage study teams
- Public and private team visibility
- Team member management (add/remove users)
- Team search functionality (by name, prefix)
- Team join requests system

### 3. Quiz System
- Create custom quizzes with multiple question types:
  - Multiple choice questions
  - Support for multiple answers
- Associate quizzes with teams
- Quiz taking and automatic grading
- Quiz results and statistics
- Pagination support for quiz listings
- Separate endpoints for quiz creation and quiz taking

### 4. File Sharing
- Upload and manage files within teams
- File metadata tracking
- Secure file storage

### 5. Real-time Messaging
- WebSocket-based real-time messaging
- Direct messages between users
- Team-based group messaging
- Message history and persistence

### 6. Friend Requests
- Send and manage friend requests
- Accept/reject friend connections
- Friend list management

### 7. Voice & Video Rooms
- Voice room creation for teams
- Screen sharing capabilities
- WebRTC signaling support
- Room capacity management (max 10 users)
- Private and group call modes
- Presenter/screen sharing controls

### 8. Event Management
- Create and manage team events
- Event scheduling and tracking

## API Endpoints

### User Endpoints
- `POST /users/signup` - Register new user
- `GET /users/:id` - Get user details
- `GET /users` - List all users
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Team Endpoints
- `POST /teams` - Create team
- `GET /teams/:id` - Get team by ID
- `GET /teams` - List all teams
- `GET /teams/search` - Search teams by prefix
- `GET /teams/by-name` - Get team by exact name
- `PUT /teams/:id` - Update team
- `DELETE /teams/:id` - Delete team
- `POST /teams/addUserToTeam` - Add user to team
- `DELETE /teams/deleteUserFromTeam` - Remove user from team

### Quiz Endpoints (Protected)
- `POST /quizzes` - Create quiz
- `GET /quizzes/:id` - Get quiz with answers
- `GET /quizzes/:id/test` - Get quiz for taking (no answers)
- `POST /quizzes/:id/test` - Submit quiz answers
- `GET /quizzes/user/:userId` - List user's quizzes
- `GET /quizzes/team/:teamId` - List team quizzes

### Real-time Communication
- `GET /messages/connect?token=<JWT>` - WebSocket connection for messaging

### Additional Features
- File upload/download endpoints
- Friend request management endpoints
- Team request endpoints
- Event management endpoints
- Voice room WebSocket endpoints

## Security Features

### Authentication & Authorization
- JWT-based token authentication
- Bearer token validation middleware
- Secure password storage using bcrypt
- Token-based WebSocket authentication

### Security Measures
- CORS configuration with allowed origins
- Input validation on all endpoints
- Protected routes requiring authentication
- Firebase security rules integration

## Real-time Features

### WebSocket Hub Architecture
The application implements a custom WebSocket hub for real-time communications:
- **Hub**: Central message broker managing all connections
- **Client**: Individual WebSocket client connections
- **Message Types**: Direct messages and team messages
- **Broadcast System**: Efficient message distribution to relevant users

### Voice Rooms
- WebRTC signaling through WebSocket
- Room state management
- Screen sharing coordination
- User presence tracking

## Development & Deployment

### Running the Server
```bash
go mod tidy           # Install dependencies
go run main.go        # Start server on port 8080
```

### Environment Configuration
Required environment variables:
- `FIREBASE_DATABASE_URL` - Firebase Realtime Database URL
- `FIREBASE_CREDENTIALS_PATH` - Path to Firebase Admin SDK credentials
- `JWT_SECRET` - Secret key for JWT signing
- `GIN_MODE` - Server mode (debug/release)

### API Documentation
- **Swagger UI**: Available at `http://localhost:8080/swagger/index.html`
- Auto-generated from code annotations
- Interactive API testing interface

### Testing
- Unit tests for controllers and services
- Test utilities and mock data
- Test files located in `/tests` directory

## Code Quality & Standards

### Code Statistics
- **92 Go source files**
- **~4,842 lines** of code in core modules (controllers, services, models, routes)
- Well-structured modular architecture

### Design Patterns
- **Repository Pattern**: Abstraction over database operations
- **Service Layer Pattern**: Business logic separation
- **Dependency Injection**: Services injected into controllers
- **Middleware Pattern**: Authentication and logging
- **DTO Pattern**: Separate data transfer objects

## Use Cases

This backend API is designed to support:
1. **Study Groups**: Students forming teams to collaborate on coursework
2. **Knowledge Sharing**: Quiz creation for peer learning and assessment
3. **Real-time Collaboration**: Live messaging and voice calls for study sessions
4. **Resource Management**: Centralized file sharing for study materials
5. **Social Learning**: Friend connections and team networking
6. **Virtual Study Rooms**: Voice/video rooms with screen sharing for group study

## Future Considerations

Based on the codebase structure, potential areas for enhancement:
- Advanced analytics and statistics
- More quiz question types
- Enhanced file management with version control
- Video recording capabilities
- Mobile push notifications
- Calendar integration for events
- Advanced search and filtering

## Conclusion

ProiectColectivBackEnd is a comprehensive, well-architected backend system for collaborative learning. It combines modern Go web development practices with real-time communication features, providing a solid foundation for a full-featured study collaboration platform. The clean separation of concerns, robust security features, and extensive API coverage make it a production-ready solution for educational team collaboration.
