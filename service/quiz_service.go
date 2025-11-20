package service

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/SerbanEduard/ProiectColectivBackEnd/mappers"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/utils"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
	ErrForbidden        = errors.New("forbidden")
)

const (
	teamNotFound  = "team not found"
	userNotFound  = "user not found"
	userNotInTeam = "user not in team"
	quizNotFound  = "quiz not found"
	NotFoundError = "not found"
)

type QuizServiceInterface interface {
	CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error)
	GetQuizWithAnswersById(id string) (entity.Quiz, error)
	GetQuizWithoutAnswersById(id string) (dto.ReadQuizResponse, error)
	SolveQuiz(request dto.SolveQuizRequest, userId string) (dto.SolveQuizResponse, error)
	GetQuizzesByUser(userId string, pageSize int, lastKey string) ([]dto.ReadQuizResponse, string, error)
	GetQuizzesByTeam(userId string, teamId string, pageSize int, lastKey string) ([]dto.ReadQuizResponse, string, error)
}

type QuizService struct {
	teamRepo TeamRepositoryInterface
	userRepo UserRepositoryInterface
	quizRepo persistence.QuizRepositoryInterface
}

func NewQuizService() *QuizService {
	return &QuizService{
		teamRepo: persistence.NewTeamRepository(),
		userRepo: persistence.NewUserRepository(),
		quizRepo: persistence.NewQuizRepository(),
	}
}

func NewQuizServiceWithRepo(teamRepo TeamRepositoryInterface, userRepo UserRepositoryInterface, quizRepo persistence.QuizRepositoryInterface) *QuizService {
	return &QuizService{
		teamRepo: teamRepo,
		userRepo: userRepo,
		quizRepo: quizRepo,
	}
}

func (qs *QuizService) isUserInTeam(userId string, teamId string) (bool, error) {
	user, err := qs.userRepo.GetByID(userId)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return false, fmt.Errorf("%w: %s", ErrResourceNotFound, userNotFound)
		}
		return false, err
	}

	ok := 0

	if user.TeamsIds == nil {
		return false, fmt.Errorf("%w: %s", ErrForbidden, teamNotFound)
	}

	teamIds := *user.TeamsIds

	for _, team := range teamIds {
		if team == teamId {
			ok = 1
		}
	}

	if ok == 0 {
		return false, fmt.Errorf("%w: %s", ErrForbidden, userNotInTeam)
	}

	return true, nil
}

func (qs *QuizService) CreateQuiz(request entity.Quiz) (dto.CreateQuizResponse, error) {
	if err := validator.ValidateCreateQuizRequest(request); err != nil {
		return dto.CreateQuizResponse{}, err
	}

	team, err := qs.teamRepo.GetTeamById(request.TeamID)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return dto.CreateQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, teamNotFound)
		}
		return dto.CreateQuizResponse{}, err
	}

	if isPartOf, err := qs.isUserInTeam(request.UserID, team.Id); isPartOf == false {
		return dto.CreateQuizResponse{}, err
	}

	for i := range request.Questions {
		questionID, err := utils.GenerateID()
		if err != nil {
			return dto.CreateQuizResponse{}, err
		}
		request.Questions[i].ID = questionID
	}

	id, err := utils.GenerateID()
	if err != nil {
		return dto.CreateQuizResponse{}, err
	}
	request.ID = id

	err = qs.quizRepo.Create(request)
	if err != nil {
		return dto.CreateQuizResponse{}, err
	}

	return dto.NewCreateQuizResponse(id), nil
}

func (qs *QuizService) GetQuizWithAnswersById(id string) (entity.Quiz, error) {
	if err := validator.ValidateQuizId(id); err != nil {
		return entity.Quiz{}, err
	}

	quiz, err := qs.quizRepo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return entity.Quiz{}, fmt.Errorf("%w: %s", ErrResourceNotFound, quizNotFound)
		}
		return entity.Quiz{}, err
	}

	return quiz, nil
}

func (qs *QuizService) GetQuizWithoutAnswersById(id string) (dto.ReadQuizResponse, error) {
	if err := validator.ValidateQuizId(id); err != nil {
		return dto.ReadQuizResponse{}, err
	}

	quiz, err := qs.quizRepo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return dto.ReadQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, quizNotFound)
		}
		return dto.ReadQuizResponse{}, err
	}

	quizWithoutAnswers := mappers.MapDomainToReadDTO(quiz)
	return quizWithoutAnswers, nil
}

func (qs *QuizService) SolveQuiz(request dto.SolveQuizRequest, userId string) (dto.SolveQuizResponse, error) {
	if err := validator.ValidateQuizId(request.QuizID); err != nil {
		return dto.SolveQuizResponse{}, err
	}

	questionsSubmitted := request.Attempts
	sort.Slice(questionsSubmitted, func(i, j int) bool {
		return questionsSubmitted[i].QuestionID < questionsSubmitted[j].QuestionID
	})
	quiz, err := qs.quizRepo.GetById(request.QuizID)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return dto.SolveQuizResponse{}, fmt.Errorf("%w: %s", ErrResourceNotFound, quizNotFound)
		}
		return dto.SolveQuizResponse{}, err
	}

	if isPartOf, err := qs.isUserInTeam(userId, quiz.TeamID); isPartOf == false {
		return dto.SolveQuizResponse{}, err
	}

	questions := quiz.Questions
	sort.Slice(questions, func(i, j int) bool {
		return questions[i].ID < questions[j].ID
	})

	// Validate the solve quiz request with questions
	if err := validator.ValidateSolveQuizRequest(request, questions); err != nil {
		return dto.SolveQuizResponse{}, err
	}

	for i, question := range questions {
		submitted := questionsSubmitted[i]
		if err := validator.ValidateQuestionSubmission(submitted, question); err != nil {
			return dto.SolveQuizResponse{}, err
		}
	}

	allCorrect := true
	questionResponses := make([]dto.SolveQuestionResponse, len(questions))

	for i, question := range questions {
		submitted := questionsSubmitted[i]
		correctFields := question.Answers
		sort.Slice(correctFields, func(i, j int) bool { return correctFields[i] < correctFields[j] })
		submittedFields := submitted.Answer
		sort.Slice(submittedFields, func(i, j int) bool { return submittedFields[i] < submittedFields[j] })
		isCorrect := slices.Equal(correctFields, submittedFields)
		if !isCorrect {
			allCorrect = false
		}
		questionResponses[i] = dto.NewSolveQuestionResponse(question.ID, isCorrect, correctFields)
	}

	return dto.SolveQuizResponse{
		IsCorrect:         allCorrect,
		QuestionResponses: questionResponses,
	}, nil
}

func (qs *QuizService) GetQuizzesByUser(userId string, pageSize int, lastKey string) ([]dto.ReadQuizResponse, string, error) {
	if err := validator.ValidateGetQuizzesByUserRequest(userId, pageSize); err != nil {
		return nil, "", err
	}

	_, err := qs.userRepo.GetByID(userId)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return nil, "", fmt.Errorf("%w: %s", ErrResourceNotFound, userNotFound)
		}
		return nil, "", err
	}

	quizzes, newKey, err := qs.quizRepo.GetByUser(userId, pageSize, lastKey)
	if err != nil {
		return nil, "", err
	}
	results := make([]dto.ReadQuizResponse, 0, len(quizzes))
	for _, quiz := range quizzes {
		if _, err := qs.isUserInTeam(userId, quiz.TeamID); err != nil {
			return nil, "", err
		}
		quizDTO := mappers.MapDomainToReadDTO(quiz)
		results = append(results, quizDTO)
	}

	return results, newKey, nil
}

func (qs *QuizService) GetQuizzesByTeam(userId string, teamId string, pageSize int, lastKey string) ([]dto.ReadQuizResponse, string, error) {
	if err := validator.ValidateGetQuizzesByTeamRequest(userId, teamId, pageSize); err != nil {
		return nil, "", err
	}

	_, err := qs.teamRepo.GetTeamById(teamId)
	if err != nil {
		if strings.Contains(err.Error(), NotFoundError) {
			return nil, "", fmt.Errorf("%w: %s", ErrResourceNotFound, teamNotFound)
		}
		return nil, "", err
	}

	if isPartOf, err := qs.isUserInTeam(userId, teamId); isPartOf == false {
		return nil, "", err
	}

	quizzes, newKey, err := qs.quizRepo.GetByTeam(teamId, pageSize, lastKey)
	if err != nil {
		return nil, "", err
	}
	results := make([]dto.ReadQuizResponse, 0, len(quizzes))
	for _, quiz := range quizzes {
		quizDTO := mappers.MapDomainToReadDTO(quiz)
		results = append(results, quizDTO)
	}

	return results, newKey, nil
}
