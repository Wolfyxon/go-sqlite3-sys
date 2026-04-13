package sqlite

import (
	"database/sql"
	"database/sql/driver"
)

/*
#cgo LDFLAGS: -lsqlite3
*/
import "C"

type sqliteDriver struct{}

var _ = (*sqliteDriver)(nil)         // allows to import the module without using any stuff declared here
var DriverName string = "sqlite-sys" // same reason but hey you can import the name

func init() {
	sql.Register(DriverName, &sqliteDriver{})
}

func (d *sqliteDriver) Open(filePath string) (driver.Conn, error) {
	conn := &sqliteConn{
		filePath: filePath,
	}

	_, err := conn.Begin()

	return conn, err
}
