package sqlite

/*
#include <string.h>
#include <sqlite3.h>
*/
import "C"
import (
	"database/sql/driver"
	"fmt"
	"io"
	"unsafe"
)

type valueType int

const (
	valueTypeUnknown = iota
	valueTypeInt
	valueTypeFloat
	valueTypeText
	valueTypeBlob
	valueTypeNull
)

type sqliteRows struct {
	stmt *sqliteStmt
}

func (r *sqliteRows) Columns() []string {
	res := []string{}

	for i := range r.getColumnCount() {
		res = append(res, r.getColumnName(i))
	}

	return res
}

func (r *sqliteRows) Next(dest []driver.Value) error {
	stepRes := C.sqlite3_step(r.stmt.handle)

	if stepRes == C.SQLITE_ROW {
		r.processCurrentRow(dest)
		return nil
	}

	if stepRes == C.SQLITE_DONE {
		return io.EOF
	}

	if stepRes == C.SQLITE_ERROR {
		return fmt.Errorf("SQLite error")
	}

	if stepRes == C.SQLITE_BUSY {
		return fmt.Errorf("Failed to lock database")
	}

	return nil
}

func (r *sqliteRows) Close() error {
	res := C.sqlite3_reset(r.stmt.handle)

	if res != 0 {
		return fmt.Errorf("Failed to reset SQLite statement: %d", res)
	}

	return nil
}

func (r *sqliteRows) processCurrentRow(dest []driver.Value) {
	for i := range r.getColumnCount() {
		valType := r.getColumnValueType(i)
		cI := C.int(i)
		handle := r.stmt.handle

		switch valType {
		case valueTypeInt:
			dest[i] = int(C.sqlite3_column_int(handle, cI))
		case valueTypeFloat:
			dest[i] = float64(C.sqlite3_column_double(handle, cI))
		case valueTypeText:
			uCharPtr := C.sqlite3_column_text(handle, cI)
			dest[i] = cUcharStrToCharStr(uCharPtr)
		case valueTypeBlob:
			cSize := C.sqlite3_column_bytes(handle, cI)
			cBlob := C.sqlite3_column_blob(handle, cI)

			buf := make([]byte, int(cSize))
			C.memcpy(cBlob, unsafe.Pointer(&buf[0]), C.size_t(cSize))

		case valueTypeNull:
			dest[i] = nil
		}
	}
}

func (r *sqliteRows) getColumnValueType(i int) valueType {
	cType := C.sqlite3_column_type(r.stmt.handle, C.int(i))

	return valueType(cType)
}

func (r *sqliteRows) getColumnName(i int) string {
	cName := C.sqlite3_column_name(r.stmt.handle, C.int(i))

	if cName == nil {
		return ""
	}

	return C.GoString(cName)
}

func (r *sqliteRows) getColumnCount() int {
	cCount := C.sqlite3_column_count(r.stmt.handle)
	return int(cCount)
}
