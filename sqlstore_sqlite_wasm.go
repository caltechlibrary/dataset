//go:build wasip1

package dataset

import (
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const Sqlite3DriverName = "sqlite3"
