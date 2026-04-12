package sqlite

/*
#include <sqlite3.h>
*/
import "C"

// TODO: Make this thread safe

type sqliteResult struct {
	stmt *sqliteStmt
}

func (r *sqliteResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r *sqliteResult) RowsAffected() (int64, error) {
	cChanges := C.sqlite3_changes64(r.stmt.connection.handle)
	return int64(cChanges), nil
}
