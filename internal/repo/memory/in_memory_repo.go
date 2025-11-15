package memory

import (
	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type InMemoryRepository struct {
	prs   map[string]*api.PullRequest
	users map[string]*api.User
	teams map[string]*api.Team
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		prs:   make(map[string]*api.PullRequest),
		users: make(map[string]*api.User),
		teams: make(map[string]*api.Team),
	}
}

func (r *InMemoryRepository) CreateUser(user *api.User) {
	r.users[user.UserID] = user
}

func (r *InMemoryRepository) GetUser(id string) (*api.User, bool) {
	user, ok := r.users[id]
	return user, ok
}

func (r *InMemoryRepository) UpdateUser(user *api.User) (*api.User, bool) {
	_, ok := r.users[user.UserID]
	if !ok {
		return nil, false
	}
	r.users[user.UserID] = user

	if team, ok := r.teams[user.TeamName]; ok {
		for i, member := range team.Members {
			if member.UserID == user.UserID {
				team.Members[i].IsActive = user.IsActive
			}
		}
		r.teams[user.TeamName] = team
	}

	return user, ok
}

func (r *InMemoryRepository) CreatePullRequest(pr *api.PullRequest) {
	r.prs[pr.PullRequestID] = pr
}

func (r *InMemoryRepository) CreateTeam(team *api.Team) {
	r.teams[team.TeamName] = team
}

func (r *InMemoryRepository) UpdateTeam(team *api.Team) (*api.Team, bool) {
	_, ok := r.teams[team.TeamName]
	if !ok {
		return nil, false
	}

	r.teams[team.TeamName] = team
	return team, true
}

func (r *InMemoryRepository) GetPullRequest(id string) (*api.PullRequest, bool) {
	pr, ok := r.prs[id]
	return pr, ok
}

func (r *InMemoryRepository) GetPullRequestsByUser(userID string) []api.PullRequestShort {
	result := make([]api.PullRequestShort, 0)
	for _, pr := range r.prs {
		for _, reviewerID := range pr.AssignedReviewers {
			if reviewerID == userID {
				shortPr := api.PullRequestShort{
					PullRequestID:   pr.PullRequestID,
					PullRequestName: pr.PullRequestName,
					AuthorID:        pr.AuthorID,
					Status:          api.PullRequestShortStatus(pr.Status),
				}
				result = append(result, shortPr)
				break
			}
		}
	}
	return result
}

func (r *InMemoryRepository) GetTeamByName(name string) (*api.Team, bool) {
	v, ok := r.teams[name]
	return v, ok
}

func (r *InMemoryRepository) CheckUserExists(id string) bool {
	_, ok := r.users[id]
	return ok
}
func (r *InMemoryRepository) CheckTeamExists(name string) bool {
	_, ok := r.teams[name]
	return ok
}

func (r *InMemoryRepository) PickRandomReviewers(authorID string, n int) []string {
	result := make([]string, 0, n)
	for userID, user := range r.users {
		if userID != authorID {
			result = append(result, user.UserID)
		}
		if len(result) == n {
			break
		}
	}
	return result
}

func (r *InMemoryRepository) UpdatePullRequest(pr *api.PullRequest) (*api.PullRequest, bool) {
	_, ok := r.prs[pr.PullRequestID]
	if !ok {
		return nil, false
	}
	r.prs[pr.PullRequestID] = pr
	return pr, true
}

func InitRepo() *InMemoryRepository {
	repo := NewInMemoryRepository()

	backend := &api.Team{TeamName: "backend"}
	frontend := &api.Team{TeamName: "frontend"}

	users := []*api.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{UserID: "u3", Username: "Charlie", TeamName: "backend", IsActive: false},
		{UserID: "u4", Username: "Dana", TeamName: "frontend", IsActive: true},
		{UserID: "u5", Username: "Eugene", TeamName: "qa", IsActive: true},
	}

	for _, u := range users {
		repo.CreateUser(u)
	}

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
