package user

import (
	"database/sql"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/rs/zerolog/log"
)

const (
	userTableName      = "users"
	profilesTableName  = "user_profiles"
	followersTableName = "user_followers"
)

type repo struct {
	db *sql.DB
}

func newUserPostgresRepository(d *sql.DB) Repository {
	return &repo{
		db: d,
	}
}

func (r *repo) Get(id string) (*entity.User, error) {
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE id = $1", userTableName)
	log.Debug().Msg("Query: " + query)
	row := r.db.QueryRow(query, id)
	user := entity.User{}

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) GetByUsername(username string) (*entity.User, error) {
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE username = $1", userTableName)
	log.Debug().Msg("Query: " + query)
	row := r.db.QueryRow(query, username)
	user := entity.User{}

	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repo) Store(user *entity.User) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) VALUES ($1, $2) RETURNING id", userTableName)
	log.Debug().Msg("Query: " + query)
	row := r.db.QueryRow(query, user.Username, user.PasswordHash)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *repo) Update(user *entity.User) error {
	query := fmt.Sprintf("UPDATE %s SET username = $1, password_hash = $2 WHERE id = $3", userTableName)
	log.Debug().Msg("Query: " + query)
	_, err := r.db.Exec(query, user.Username, user.PasswordHash, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Remove(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", userTableName)
	log.Debug().Msg("Query: " + query)
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddFollower(userID, followerID string) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, follower_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", followersTableName)
	log.Debug().Msg("Query: " + query)
	_, err := r.db.Exec(query, userID, followerID)
	if err != nil {
		return err
	}

	return nil
}
