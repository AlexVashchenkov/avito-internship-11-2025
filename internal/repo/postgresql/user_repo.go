package postgresql

import (
	"database/sql"
	"errors"
	"log/slog"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *api.User) {
	_, err := r.db.Exec(`
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
	`, user.UserID, user.Username, user.TeamName, user.IsActive)

	if err != nil {
		slog.Error("failed to create user",
			"user_id", user.UserID,
			"team_name", user.TeamName,
			"error", err.Error(),
		)
	} else {
		slog.Info("user created",
			"user_id", user.UserID,
		)
	}
}

func (r *UserRepository) GetUser(id string) (*api.User, bool) {
	row := r.db.QueryRow(`
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`, id)

	u := &api.User{}
	err := row.Scan(&u.UserID, &u.Username, &u.TeamName, &u.IsActive)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("user not found", "user_id", id)
		return nil, false
	}
	if err != nil {
		slog.Error("failed to load user", "user_id", id, "error", err.Error())
		return nil, false
	}

	return u, true
}

func (r *UserRepository) UpdateUser(u *api.User) (*api.User, bool) {
	_, err := r.db.Exec(`
		UPDATE users
		SET username=$2, team_name=$3, is_active=$4
		WHERE user_id=$1
	`, u.UserID, u.Username, u.TeamName, u.IsActive)

	if err != nil {
		slog.Error("failed to update user",
			"user_id", u.UserID,
			"error", err.Error(),
		)
		return u, false
	}

	slog.Info("user updated", "user_id", u.UserID)
	return u, true
}

func (r *UserRepository) PickRandomReviewers(authorID, teamName string, n int) []string {
	rows, err := r.db.Query(`
		SELECT user_id
		FROM users
		WHERE user_id <> $1
		  AND team_name=$2
		  AND is_active = true
		ORDER BY random()
		LIMIT $3
	`, authorID, teamName, n)

	if err != nil {
		slog.Error("failed to pick random reviewers",
			"author_id", authorID,
			"team_name", teamName,
			"error", err.Error(),
		)
		return nil
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			slog.Error("failed to scan reviewer row", "error", err.Error())
			continue
		}
		res = append(res, id)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error iterating reviewer rows",
			"author_id", authorID,
			"team_name", teamName,
			"error", err.Error(),
		)
		return nil
	}

	slog.Info("reviewers selected",
		"author_id", authorID,
		"team_name", teamName,
		"count", len(res),
	)

	return res
}
