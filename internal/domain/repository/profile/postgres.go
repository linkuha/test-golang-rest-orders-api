package profile

import (
	"database/sql"
	"fmt"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
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

func (r *repo) GetByUserID(userID string) (*entity.Profile, error) {
	query := fmt.Sprintf(`SELECT first_name, last_name, middle_name, 
		TRIM(CONCAT_WS(' ', last_name, first_name, middle_name)) AS full_name, sex, age FROM %s WHERE user_id = ?`, tableName)
	row := r.db.QueryRow(query, userID)
	profile := entity.Profile{}

	err := row.Scan(&profile.FirstName, &profile.LastName, &profile.MiddleName, &profile.FullName, &profile.Sex, &profile.Age)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *repo) Store(profile *entity.Profile) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (user_id, first_name, last_name, middle_name, sex, age) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`, tableName)
	row := r.db.QueryRow(query, profile.UserID, profile.FirstName, profile.LastName, profile.MiddleName, profile.Sex, profile.Age)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Update(profile *entity.Profile) error {
	query := fmt.Sprintf(`UPDATE %s SET first_name = ?, last_name = ?, middle_name = ?, sex = ?, age = ? WHERE user_id = ?`, tableName)
	_, err := r.db.Exec(query, profile.FirstName, profile.LastName, profile.MiddleName, profile.Sex, profile.Age, profile.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) RemoveByUserID(userID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = ?`, tableName)
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
