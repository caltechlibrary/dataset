package dataset

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
)

type DSQuery struct {
	CName      string   `json:"c_name,omitempty"`
	Stmt       string   `json:"stmt,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
	ds         *Collection
	db         *sql.DB
}

func (app *DSQuery) RunQuery(in io.Reader, out io.Writer, eout io.Writer, cName string, stmt string, jsonl bool, params []string) error {
	app.CName = cName
	app.Stmt = stmt
	ds, err := Open(cName)
	if err != nil {
		return err
	}
	defer ds.Close()
	app.ds = ds
	app.db = ds.SQLStore.db

	var rows *sql.Rows
	if len(params) > 0 {
		args := []interface{}{}
		for _, val := range params {
			args = append(args, val)
		}
		rows, err = app.db.Query(app.Stmt, args...)
	} else {
		rows, err = app.db.Query(app.Stmt)
	}
	if err != nil {
		return fmt.Errorf("stmt: %s, %s", app.Stmt, err)
	}
	src := []byte{}
	if (! jsonl) {
		src = append(src, '[')
	}
	i := 0
	for rows.Next() {
		// Get our row values
		objSrc := []byte{}
		if err := rows.Scan(&objSrc); err != nil {
			return err
		}
		// flatten out objSrc if needed
		if jsonl && bytes.Contains(objSrc, []byte{'\n'}) {
			objSrc = FmtJSONL(objSrc)
		}
		if i > 0 {
			if jsonl {
				src = append(src, '\n')
			} else {
				src = append(src, ',')
			}
		}
		// Object should be formatted as one line.
		src = append(src, objSrc...)
		i++
	}
	if (!jsonl) {
		src = append(src, ']')
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}
