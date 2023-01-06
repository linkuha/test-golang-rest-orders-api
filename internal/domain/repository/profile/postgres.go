package profile

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/rs/zerolog/log"
)

const tableName = "user_profiles"

type repo struct {
	db *sql.DB
}

func newProfilePostgresRepository(d *sql.DB) Repository {
	return &repo{
		db: d,
	}
}

func (r *repo) GetByUserID(ctx context.Context, userID string) (*entity.Profile, error) {
	query := fmt.Sprintf(`SELECT user_id, first_name, last_name, middle_name,
		TRIM(CONCAT_WS(' ', last_name, first_name, middle_name)) AS full_name, sex, age FROM %s WHERE user_id = $1`, tableName)
	log.Debug().Msg("Query: " + query)

	row := r.db.QueryRowContext(ctx, query, userID)
	profile := entity.Profile{}

	err := row.Scan(&profile.UserID, &profile.FirstName, &profile.LastName, &profile.MiddleName, &profile.FullName, &profile.Sex, &profile.Age)
	if err != nil {
		return nil, errs.HandleErrorDB(err)
	}
	return &profile, nil
}

func (r *repo) Store(ctx context.Context, profile *entity.Profile) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (user_id, first_name, last_name, middle_name, sex, age) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, tableName)
	log.Debug().Msg("Query: " + query)

	row := r.db.QueryRowContext(ctx, query, profile.UserID, profile.FirstName, profile.LastName, profile.MiddleName, profile.Sex, profile.Age)
	if err := row.Scan(&id); err != nil {
		return 0, errs.HandleErrorDB(err)
	}

	return id, nil
}

func (r *repo) Update(ctx context.Context, profile *entity.Profile) error {
	query := fmt.Sprintf(`UPDATE %s SET first_name = $1, last_name = $2, middle_name = $3, sex = $4, age = $5 WHERE user_id = $6`, tableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, profile.FirstName, profile.LastName, profile.MiddleName, profile.Sex, profile.Age, profile.UserID)
	if err != nil {
		return errs.HandleErrorDB(err)
	}

	return nil
}

func (r *repo) RemoveByUserID(ctx context.Context, userID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1`, tableName)
	log.Debug().Msg("Query: " + query)

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return errs.HandleErrorDB(err)
	}
	return nil
}
