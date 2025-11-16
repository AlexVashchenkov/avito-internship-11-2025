package contracts

import api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"

type TeamRepository interface {
	CreateTeam(team *api.Team)
	GetTeamByName(name string) (*api.Team, bool)
	UpdateTeam(team *api.Team) (*api.Team, bool)
	UpdateTeamMember(user *api.User) (*api.User, bool)
}
