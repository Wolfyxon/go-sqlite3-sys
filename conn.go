package sqlite

import (
	"database/sql/driver"
	"fmt"
)

/*
#include <sqlite3.h>
*/
import "C"

type sqliteConn struct {
	filePath string
	handle   *C.sqlite3
}

func (c *sqliteConn) Prepare(query string) (driver.Stmt, error) {
	if c.handle == nil {
		return nil, fmt.Errorf("Database handle is nil")
	}

	cQuery := C.CString(query)
	cLen := C.int(len(query))

	var cStmt *C.sqlite3_stmt
	var cUnusedSql *C.char

	res := C.sqlite3_prepare_v2(c.handle, cQuery, cLen, &cStmt, &cUnusedSql)

	if res != 0 {
		return nil, fmt.Errorf("sqlite error: %d", res)
	}

	return &sqliteStmt{
		handle:      cStmt,
		connection:  c,
		queryString: query,
	}, nil
}

func (c *sqliteConn) Begin() (driver.Tx, error) {

	var handle *C.sqlite3
	res := C.sqlite3_open(C.CString(c.filePath), &handle)

	if res != 0 {
		return nil, fmt.Errorf("Failed to open sqlite database: %d", res)
	}

	c.handle = handle

	return nil, nil
}

func (c *sqliteConn) Close() error {
	if c.handle == nil {
		return nil
	}

	res := C.sqlite3_close(c.handle)
	c.handle = nil

	if res != 0 {
		return fmt.Errorf("Failed to close sqlite3: %d", res)
	}

	return nil
}
