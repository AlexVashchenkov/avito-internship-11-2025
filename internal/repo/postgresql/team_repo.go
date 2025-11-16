package postgresql

import (
	"database/sql"
	"log/slog"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(t *api.Team) {
	_, err := r.db.Exec(`
		INSERT INTO teams (team_name)
		VALUES ($1)
	`, t.TeamName)

	if err != nil {
		slog.Error("failed to create team", "team_name", t.TeamName, "error", err.Error())
	} else {
		slog.Info("team created", "team_name", t.TeamName)
	}
}

func (r *TeamRepository) CreateTeamWithMembers(t *api.Team) error {
	tx, err := r.db.Begin()
	if err != nil {
		slog.Error("failed to start transaction", "error", err.Error())
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO teams (team_name)
		VALUES ($1)
	`, t.TeamName)
	if err != nil {
		slog.Error("failed to insert team", "team_name", t.TeamName, "error", err.Error())
		return err
	}

	for _, member := range t.Members {
		_, err = tx.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) 
			DO UPDATE SET 
				username = EXCLUDED.username,
				team_name = EXCLUDED.team_name,
				is_active = EXCLUDED.is_active
		`, member.UserID, member.Username, t.TeamName, member.IsActive)

		if err != nil {
			slog.Error("failed to insert team member",
				"team_name", t.TeamName,
				"user_id", member.UserID,
				"error", err.Error(),
			)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit team create", "team_name", t.TeamName, "error", err.Error())
		return err
	}

	slog.Info("team with members created", "team_name", t.TeamName)
	return nil
}

func (r *TeamRepository) GetTeamByName(name string) (*api.Team, bool) {
	team := &api.Team{TeamName: name}

	rows, err := r.db.Query(`
		SELECT user_id, username, is_active, team_name
		FROM users
		WHERE team_name=$1
	`, name)

	if err != nil {
		slog.Error("failed to load team members", "team_name", name, "error", err.Error())
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var u api.User
		if scanErr := rows.Scan(&u.UserID, &u.Username, &u.IsActive, &u.TeamName); scanErr != nil {
			slog.Error("failed to scan team member", "team_name", name, "error", scanErr.Error())
			continue
		}

		team.Members = append(team.Members, api.TeamMember{
			UserID: u.UserID, Username: u.Username, IsActive: u.IsActive,
		})
	}

	if err := rows.Err(); err != nil {
		slog.Error("error iterating team member rows", "team_name", name, "error", err.Error())
		return nil, false
	}

	return team, true
}

func (r *TeamRepository) UpdateTeam(t *api.Team) (*api.Team, bool) {
	_, err := r.db.Exec(`
		UPDATE teams SET team_name=$1 WHERE team_name=$1
	`, t.TeamName)

	if err != nil {
		slog.Error("failed to update team", "team_name", t.TeamName, "error", err.Error())
		return t, false
	}

	slog.Info("team updated", "team_name", t.TeamName)
	return t, true
}

func (r *TeamRepository) UpdateTeamMember(u *api.User) (*api.User, bool) {
	_, err := r.db.Exec(`
		UPDATE users
		SET username=$2, team_name=$3, is_active=$4
		WHERE user_id=$1
	`, u.UserID, u.Username, u.TeamName, u.IsActive)

	if err != nil {
		slog.Error("failed to update team member",
			"user_id", u.UserID,
			"team_name", u.TeamName,
			"error", err.Error(),
		)
		return u, false
	}

	slog.Info("team member updated", "user_id", u.UserID)
	return u, true
}
