package service

import (
	"errors"

	"github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"
	"github.com/SerbanEduard/ProiectColectivBackEnd/model/entity"
	"github.com/SerbanEduard/ProiectColectivBackEnd/persistence"
	"github.com/SerbanEduard/ProiectColectivBackEnd/validator"
)

type TeamRequestRepositoryInterface interface {
	Create(req *entity.TeamRequest) error
	GetById(id string) (*entity.TeamRequest, error)
	GetAll() ([]*entity.TeamRequest, error)
	GetByUserId(userId string) ([]*entity.TeamRequest, error)
	Delete(id string) error
}

type TeamRequestService struct {
	teamRequestRepository TeamRequestRepositoryInterface
	userRepository        UserRepositoryInterface
	teamRepository        TeamRepositoryInterface
	teamService           TeamServiceInterface
}

func NewTeamRequestService() *TeamRequestService {
	return &TeamRequestService{
		teamRequestRepository: persistence.NewTeamRequestRepository(),
		userRepository:        persistence.NewUserRepository(),
		teamRepository:        persistence.NewTeamRepository(),
		teamService:           NewTeamService(),
	}
}

func NewTeamRequestServiceWithRepo(
	trRepo TeamRequestRepositoryInterface,
	userRepo UserRepositoryInterface,
	teamRepo TeamRepositoryInterface,
	teamService TeamServiceInterface,
) *TeamRequestService {
	return &TeamRequestService{
		teamRequestRepository: trRepo,
		userRepository:        userRepo,
		teamRepository:        teamRepo,
		teamService:           teamService,
	}
}

type TeamServiceInterface interface {
	AddUserToTeam(idUser string, idTeam string) (*entity.User, *entity.Team, error)
	DeleteUserFromTeam(idUser string, idTeam string) (*entity.User, *entity.Team, error)
	GetTeamById(id string) (*entity.Team, error)
	GetXTeamsByPrefix(prefix string, x int) ([]*entity.Team, error)
	GetTeamsByName(name string) ([]*entity.Team, error)
	GetAll() ([]*entity.Team, error)
	Update(team *entity.Team) error
	Delete(id string) error
}

func (trs *TeamRequestService) CreateTeamRequest(req *dto.TeamRequestCreateDTO) (*entity.TeamRequest, error) {
	if err := validator.ValidateTeamRequestCreateDTO(req); err != nil {
		return nil, err
	}

	user, err := trs.userRepository.GetByID(req.UserID)
	if err != nil {
		return nil, errors.New("user does not exist")
	}

	team, err := trs.teamRepository.GetTeamById(req.TeamID)
	if err != nil {
		return nil, errors.New("team does not exist")
	}

	for _, uid := range team.UsersIds {
		if uid == req.UserID {
			return nil, errors.New("user is already a member of this team")
		}
	}

	existingRequests, err := trs.teamRequestRepository.GetByUserId(req.UserID)
	if err == nil {
		for _, r := range existingRequests {
			if r.TeamID == req.TeamID {
				return nil, errors.New("a pending request already exists for this user and team")
			}
		}
	}

	id, err := generateID()
	if err != nil {
		return nil, err
	}

	newReq := entity.NewTeamRequest(id, user.ID, team.Id)
	if err := trs.teamRequestRepository.Create(newReq); err != nil {
		return nil, err
	}

	return newReq, nil
}

func (trs *TeamRequestService) AcceptTeamRequest(id string) (*entity.User, *entity.Team, error) {
	req, err := trs.teamRequestRepository.GetById(id)
	if err != nil {
		return nil, nil, errors.New("team request not found")
	}

	user, team, err := trs.teamService.AddUserToTeam(req.UserID, req.TeamID)
	if err != nil {
		return nil, nil, err
	}

	_ = trs.teamRequestRepository.Delete(id)
	return user, team, nil
}

func (trs *TeamRequestService) RejectTeamRequest(id string) error {
	req, err := trs.teamRequestRepository.GetById(id)
	if err != nil {
		return errors.New("team request not found")
	}
	return trs.teamRequestRepository.Delete(req.Id)
}

func (trs *TeamRequestService) GetAll() ([]*entity.TeamRequest, error) {
	return trs.teamRequestRepository.GetAll()
}

func (trs *TeamRequestService) GetByUserId(userId string) ([]*entity.TeamRequest, error) {
	return trs.teamRequestRepository.GetByUserId(userId)
}
