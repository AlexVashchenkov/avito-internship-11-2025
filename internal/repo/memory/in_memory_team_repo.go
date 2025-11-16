package memory

import (
	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type InMemoryTeamRepository struct {
	teams map[string]*api.Team
}

func NewInMemoryTeamRepository() *InMemoryTeamRepository {
	return &InMemoryTeamRepository{
		teams: make(map[string]*api.Team),
	}
}

func (r *InMemoryTeamRepository) CreateTeam(team *api.Team) {
	r.teams[team.TeamName] = team
}

func (r *InMemoryTeamRepository) UpdateTeam(team *api.Team) (*api.Team, bool) {
	_, ok := r.teams[team.TeamName]
	if !ok {
		return nil, false
	}

	r.teams[team.TeamName] = team
	return team, true
}

func (r *InMemoryTeamRepository) GetTeamByName(name string) (*api.Team, bool) {
	v, ok := r.teams[name]
	return v, ok
}

func (r *InMemoryTeamRepository) UpdateTeamMember(user *api.User) (*api.User, bool) {
	team, ok := r.teams[user.TeamName]
	if ok {
		for i, member := range team.Members {
			if member.UserID == user.UserID {
				team.Members[i].IsActive = user.IsActive
			}
		}
		r.teams[user.TeamName] = team
	}

	return user, ok
}

func InitTeamRepo() *InMemoryTeamRepository {
	repo := NewInMemoryTeamRepository()

	backend := &api.Team{TeamName: "backend"}
	frontend := &api.Team{TeamName: "frontend"}

	backend.Members = []api.TeamMember{
		{UserID: "u1", Username: "Alice", IsActive: true},
		{UserID: "u2", Username: "Bob", IsActive: true},
		{UserID: "u3", Username: "Charlie", IsActive: false},
	}
	frontend.Members = []api.TeamMember{
		{UserID: "u4", Username: "Dana", IsActive: true},
	}

	repo.CreateTeam(backend)
	repo.CreateTeam(frontend)

	return repo
}
