-- +goose Up

-- Teams
INSERT INTO teams (team_name)
VALUES ('backend'), ('frontend'), ('qa')
ON CONFLICT DO NOTHING;

-- Users: backend
INSERT INTO users (user_id, username, team_name, is_active)
VALUES
    ('u1', 'Alice',   'backend', TRUE),
    ('u2', 'Bob',     'backend', TRUE),
    ('u3', 'Charlie', 'backend', FALSE)
ON CONFLICT DO NOTHING;

-- Users: frontend
INSERT INTO users (user_id, username, team_name, is_active)
VALUES
    ('u4', 'Dana', 'frontend', TRUE)
ON CONFLICT DO NOTHING;

-- Users: qa
INSERT INTO users (user_id, username, team_name, is_active)
VALUES
    ('u5', 'Eugene', 'qa', TRUE)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM users WHERE user_id IN ('u1','u2','u3','u4','u5');
DELETE FROM teams WHERE team_name IN ('backend','frontend','qa');