package validator

import "github.com/SerbanEduard/ProiectColectivBackEnd/model/dto"

func ValidateTeamRequestCreateDTO(request *dto.TeamRequestCreateDTO) error {
	validations := []func() error{
		func() error { return validateRequired(request.UserID, "userId is required") },
		func() error { return validateRequired(request.TeamID, "teamId is required") },
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}
