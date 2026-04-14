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
		return r.processCurrentRow(dest)
	}

	if stepRes == C.SQLITE_DONE {
		return io.EOF
	}

	err := getErrorFromCode(stepRes)

	if err != nil {
		C.sqlite3_reset(r.stmt.handle)
		return err
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

func (r *sqliteRows) processCurrentRow(dest []driver.Value) error {
	for i := range r.getColumnCount() {
		valType := r.getColumnValueType(i)
		cI := C.int(i)
		handle := r.stmt.handle

		switch valType {
		case valueTypeUnknown:
			return fmt.Errorf("Could not determine type of column %d", i)
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
			C.memcpy(unsafe.Pointer(&buf[0]), cBlob, C.size_t(cSize))
			dest[i] = buf
		case valueTypeNull:
			dest[i] = nil
		default:
			return fmt.Errorf("Driver bug: sqlite type ID '%d' not implemented. Column: %d", valType, i)
		}
	}

	return nil
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
