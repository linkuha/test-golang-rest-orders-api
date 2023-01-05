package errs

import (
	"database/sql"
	"database/sql/driver"
	"errors"
)

var (
	RecordNotFound = errors.New("record is not found")
)

func HandleErrorDB(e error) error {
	if e == nil {
		return nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return NewErrorWrapper(NotExist, RecordNotFound, "not found")
	}
	if errors.Is(e, sql.ErrConnDone) || errors.Is(e, driver.ErrBadConn) {
		return NewErrorWrapper(DatabaseConnection, e, "connection problem")
	}
	return NewErrorWrapper(Database, e, "db another error")
}
