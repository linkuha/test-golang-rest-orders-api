package user

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"testing"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(id, "qwerty", "password"))

	r := newUserPostgresRepository(db)
	u, err := r.Get(id)
	if err != nil {
		t.Errorf("error was not expected while get user: %s", err)
	}
	if id != u.ID {
		t.Errorf("ids is not equal")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	mock.ExpectQuery("SELECT").WithArgs(id).
		WillReturnError(fmt.Errorf("some error"))

	r := newUserPostgresRepository(db)
	u, err := r.Get(id)
	if err == nil {
		t.Errorf("was expecting an error, but there was none")
	}
	if u != nil {
		t.Errorf("was expecting nil model, but there exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "qwerty"
	id := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	mock.ExpectQuery("SELECT").WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(id, username, "password"))

	r := newUserPostgresRepository(db)
	u, err := r.GetByUsername(username)
	if err != nil {
		t.Errorf("error was not expected while get user: %s", err)
	}
	if username != u.Username {
		t.Errorf("usernames is not equal")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByUsernameError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "qwerty"

	mock.ExpectQuery("SELECT").WithArgs(username).
		WillReturnError(fmt.Errorf("some error"))

	r := newUserPostgresRepository(db)
	u, err := r.GetByUsername(username)
	if err == nil {
		t.Errorf("was expecting an error, but there was none")
	}
	if u != nil {
		t.Errorf("was expecting nil model, but there exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := entity.TestExistUser(t)

	query := fmt.Sprintf("INSERT INTO %s", userTableName)
	mock.ExpectQuery(query).WithArgs(u.Username, u.PasswordHash).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(u.ID))

	r := newUserPostgresRepository(db)
	id, err := r.Store(u)
	if err != nil {
		t.Errorf("error was not expected while insert user: %s", err)
	}
	if id == "" {
		t.Errorf("empty inserted id return")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStoreError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := entity.TestExistUser(t)

	query := fmt.Sprintf("INSERT INTO %s", userTableName)
	mock.ExpectQuery(query).WithArgs(u.Username, u.PasswordHash).
		WillReturnError(fmt.Errorf("some error"))

	r := newUserPostgresRepository(db)
	if _, err := r.Store(u); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := entity.TestExistUser(t)

	query := fmt.Sprintf("UPDATE %s", userTableName)
	mock.ExpectExec(query).WithArgs(u.Username, u.PasswordHash, u.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := newUserPostgresRepository(db)
	if err := r.Update(u); err != nil {
		t.Errorf("error was not expected while insert user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := entity.TestExistUser(t)

	query := fmt.Sprintf("UPDATE %s", userTableName)
	mock.ExpectExec(query).WithArgs(u.Username, u.PasswordHash, u.ID).
		WillReturnError(fmt.Errorf("some error"))

	r := newUserPostgresRepository(db)
	if err := r.Update(u); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRemove(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	query := fmt.Sprintf("DELETE FROM %s", userTableName)
	mock.ExpectExec(query).WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := newUserPostgresRepository(db)
	if err := r.Remove(id); err != nil {
		t.Errorf("error was not expected while insert user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestRemoveError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	id := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	query := fmt.Sprintf("DELETE FROM %s", userTableName)
	mock.ExpectExec(query).WithArgs(id).
		WillReturnError(fmt.Errorf("some error"))

	r := newUserPostgresRepository(db)
	if err := r.Remove(id); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddFollower(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u1 := entity.TestExistUser(t)
	u2 := entity.TestExistUser2(t)

	query := fmt.Sprintf("INSERT INTO %s", followersTableName)
	mock.ExpectExec(query).WithArgs(u1.ID, u2.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	r := newUserPostgresRepository(db)
	if err := r.AddFollower(u1.ID, u2.ID); err != nil {
		t.Errorf("error was not expected while insert user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestAddFollowerError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u1 := entity.TestExistUser(t)
	u2 := entity.TestExistUser2(t)

	query := fmt.Sprintf("INSERT INTO %s", followersTableName)
	mock.ExpectExec(query).WithArgs(u1.ID, u2.ID).
		WillReturnError(fmt.Errorf("some error"))

	r := newUserPostgresRepository(db)
	if err := r.AddFollower(u1.ID, u2.ID); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
