package dataset

import (
	"database/sql"
	"fmt"
	"io"
	"strings"
)

type DSQuery struct {
	CName string `json:"c_name,omitempty"`
	Stmt string `json:"stmt,omitempty"`
	Pretty bool `json:"pretty,omitempty"`
	ds *Collection
	db *sql.DB
}

func (app *DSQuery) Run(in io.Reader, out io.Writer, eout io.Writer, cName string, stmt string, params []string) error {
	app.CName = cName
	app.Stmt = stmt
	ds, err := Open(cName)
	if err != nil {
		return err
	}
	defer ds.Close()
	app.ds = ds
	if strings.Compare(ds.StoreType, SQLSTORE) != 0 {
		return fmt.Errorf("%s storage type not supported for this application", ds.StoreType)
	}
	if ds.SQLStore == nil {
		return fmt.Errorf("sqlstore failed to open")
	}
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
		return fmt.Errorf("sql error: %s", err)
	}
	src := []byte(`[`)
	i := 0
	for rows.Next() {
		// Get our row values
		obj := []byte{}
		if err := rows.Scan(&obj); err != nil  {
			return err
		}
		if i > 0 {
			src = append(src, ',')
		}
		src = append(src, obj...)
		i++
	}
	src = append(src, ']')
	err = rows.Err()
	if err != nil {
		return err
	}
	if app.Pretty {
		fmt.Fprintf(out, "%s\n", JSONIndent(src, "", "    "))
		return nil
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}
