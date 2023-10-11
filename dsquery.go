package dataset

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"path"
	"strings"
	"time"
)

type DSQuery struct {
	CName string `json:"c_name,omitempty"`
	Stmt string `json:"stmt,omitempty"`
	Pretty bool `json:"pretty,omitempty"`
	AsGrid bool `json:"as_grid,omitempty"`
	AsCSV bool `json:"csv,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
	PTIndex bool `json:"pt_index,omitempty"`
	ds *Collection
	db *sql.DB
}

// MakeGrid takes JSON source holding an array of objects and uses the attribute list to
// render a 2D grid of values where the columns match the attribute name list provided.
// If an attribute is missing a nil is inserted. MakeGrid returns the grid as JSON source
// along with an error value.
func MakeGrid(src []byte, attributes []string) ([]byte, error) {
	listOfObjects := []map[string]interface{}{}
	if err := JSONUnmarshal(src, &listOfObjects); err != nil {
		return nil, err
	}
	outerArray := []interface{}{}
	// Now that we know the column names scan the list of objects and build grid
	for _, obj := range listOfObjects {
		innerArray := []interface{}{}
		for _, attr := range attributes {
			if val, ok := obj[attr]; ok {
				innerArray = append(innerArray, val)
			} else {
				innerArray = append(innerArray, nil)
			}
		}
		outerArray = append(outerArray, innerArray)
	}
	return JSONMarshal(outerArray)
}

// MakeCSV takes JSON source holding an array of objects and uses the attribute list to
// render a CSV file from the list. It returns the CSV content as a byte slice along
// with an error.
func MakeCSV(src []byte, attributes []string) ([]byte, error) {
	listOfObjects := []map[string]interface{}{}
	if err := JSONUnmarshal(src, &listOfObjects); err != nil {
		return nil, err
	}
	// Write header row based on our attribute list.
	buf := []byte{}
	out := bytes.NewBuffer(buf)
	w := csv.NewWriter(out)
	if err := w.Write(attributes); err != nil {
		return nil, err
	}
	// Now that we know the column names scan the list of objects and build grid
	for _, obj := range listOfObjects {
		innerArray := []string{}
		for _, attr := range attributes {
			if val, ok := obj[attr]; ok {
				switch val.(type) {
				case string:
					cell := val.(string)
					innerArray = append(innerArray, cell)
				default:
					data, _ := JSONMarshal(val)
					innerArray = append(innerArray, fmt.Sprintf("%s", data))
				}
			} else {
				innerArray = append(innerArray, "")
			}
		}
		if err := w.Write(innerArray); err != nil {
			return nil, err
		}
	}
	w.Flush()
	err := w.Error()
	if err != nil {
		return nil, err
	}
	src = out.Bytes()
	return src, nil
}


// indexCollection creates a SQL replica of a dataset collection as a SQLite 3
// database called index.db inside the collection's root folder. This allows
// pairtree implementations to be queried using SQLite 3's SQL dialect.
func indexCollection(ds *Collection, index *sql.DB) error {
	// Clear the existing index if exists, create a new one.
	tName := strings.TrimSuffix(ds.Name, ".ds")
	stmt := fmt.Sprintf(`drop table if exists %s;`, tName)
	_, err := index.Exec(stmt)
	if err != nil {
		return fmt.Errorf("stmt:\n%s\n\t%s", stmt, err)
	}
	stmt = fmt.Sprintf(`create table %s (_key string, src json, updated datetime, created timestamp);`, tName)
	_, err = index.Exec(stmt)
	if err != nil {
		return fmt.Errorf("stmt:\n%s\n\t%s", stmt, err)
	}

	keys, err := ds.Keys()
	if err != nil {
		return err
	}
	tStamp := time.Now()
	stmt = fmt.Sprintf(`insert into %s (_key, src, updated, created) values (?, ?, ?, ?)`, tName)
	for _, key := range keys {
		src, err := ds.ReadJSON(key)
		if err != nil {
			return err
		}
		_, err = index.Exec(stmt, key, fmt.Sprintf("%s", src), tStamp, tStamp)
		if err != nil {
			return fmt.Errorf("stmt:\n%s\n\t$1 = %+v $2 = %+v, %s", stmt, tStamp, tStamp, err)
		}
	}
	return nil
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
	if strings.Compare(ds.StoreType, SQLSTORE) == 0 {
		if ds.SQLStore == nil {
			return fmt.Errorf("sqlstore failed to open")
		}
		app.db = ds.SQLStore.db
	} else {
		wPath := cName
		if ds.workPath != "" {
			wPath = ds.workPath
		}
		indexDSN := path.Join(wPath, "index.db")
		index, err := sql.Open("sqlite", indexDSN)
		if err != nil {
			return fmt.Errorf("failed to open index %s, %s", indexDSN, err)
		}
		defer index.Close()
		if index == nil {
			return fmt.Errorf("index failed to open")
		}
		if app.PTIndex {
			// Index the collection into a SQLite3 database. This would normally used to allow
			// dsquery to support querring a pairtree collection.
			if err = indexCollection(ds, index); err != nil {
				return fmt.Errorf("failed to index %q, %s", cName, err)
			}
		}
		app.db = index
	}
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
	if app.AsGrid {
		src, err = MakeGrid(src, app.Attributes)
		if err != nil {
			return err
		}
	}
	if app.AsCSV {
		src, err = MakeCSV(src, app.Attributes)
		if err != nil {
			return err
		}
	}
	if app.Pretty {
		fmt.Fprintf(out, "%s\n", JSONIndent(src, "", "    "))
		return nil
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}
