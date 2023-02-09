/**
 * import_export.go provides methods to import and export JSON content to
 * and from tables or CSV files.
 */
package dataset

import (
	"log"
	"os"
	"strings"
	"encoding/json"
	"encoding/csv"
	"io"
	"fmt"
	"strconv"

	// Caltech Library
	"github.com/caltechlibrary/dotpath"
)

const (
	// internal virtualize column name format string
	fmtColumnName = `column_%03d`
)


// ImportCSV takes a reader and iterates over the rows and imports them as
// a JSON records into dataset.
//BUG: returns lines processed should probably return number of rows imported
func (c *Collection) ImportCSV(buf io.Reader, idCol int, skipHeaderRow bool, overwrite bool, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		key        string
		err        error
	)
	r := csv.NewReader(buf)
	r.FieldsPerRecord = -1
	r.TrimLeadingSpace = true
	lineNo := 0
	if skipHeaderRow == true {
		lineNo++
		fieldNames, err = r.Read()
		if err != nil {
			return lineNo, fmt.Errorf("Can't read header csv table at %d, %s", lineNo, err)
		}
	}
	for {
		lineNo++
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return lineNo, fmt.Errorf("Can't read row csv table at %d, %s", lineNo, err)
		}
		var fieldName string
		record := map[string]interface{}{}
		if idCol < 0 {
			key = fmt.Sprintf("%d", lineNo)
		}
		for i, val := range row {
			if i < len(fieldNames) {
				fieldName = fieldNames[i]
				if idCol == i {
					key = val
				}
			} else {
				fieldName = fmt.Sprintf(fmtColumnName, i+1)
			}
			//Note: We need to convert the value
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				record[fieldName] = i
			} else if f, err := strconv.ParseFloat(val, 64); err == nil {
				record[fieldName] = f
			} else if strings.ToLower(val) == "true" {
				record[fieldName] = true
			} else if strings.ToLower(val) == "false" {
				record[fieldName] = false
			} else {
				val = strings.TrimSpace(val)
				if len(val) > 0 {
					record[fieldName] = val
				}
			}
		}
		if len(key) > 0 && len(record) > 0 {
			if c.HasKey(key) {
				if overwrite == true {
					err = c.Update(key, record)
					if err != nil {
						return lineNo, fmt.Errorf("can't update %+v to %s, %s", record, key, err)
					}
				} else if verboseLog {
					log.Printf("Skipping row %d, key %q, already exists", lineNo, key)
				}
			} else {
				err = c.Create(key, record)
				if err != nil {
					return lineNo, fmt.Errorf("can't create %+v to %s, %s", record, key, err)
				}
			}
		} else if verboseLog {
			log.Printf("Skipping row %d, key value missing", lineNo)
		}
		if verboseLog == true && (lineNo%1000) == 0 {
			log.Printf("%d rows processed", lineNo)
		}
	}
	return lineNo, nil
}

// ImportTable takes a [][]interface{} and iterates over the rows and
// imports them as a JSON records into dataset.
func (c *Collection) ImportTable(table [][]interface{}, idCol int, useHeaderRow bool, overwrite, verboseLog bool) (int, error) {
	var (
		fieldNames []string
		key        string
		err        error
	)
	if len(table) == 0 {
		return 0, fmt.Errorf("No data in table")
	}
	lineNo := 0
	// i.e. use the header row for field names
	if useHeaderRow == true {
		for i, val := range table[0] {
			cell, err := ValueInterfaceToString(val)
			if err == nil && strings.TrimSpace(cell) != "" {
				fieldNames = append(fieldNames, cell)
			} else {
				fieldNames = append(fieldNames, fmt.Sprintf(fmtColumnName, i))
			}
		}
		lineNo++
	}
	rowCount := len(table)
	for {
		if lineNo >= rowCount {
			break
		}
		row := table[lineNo]
		lineNo++

		var fieldName string
		record := map[string]interface{}{}
		if idCol < 0 {
			key = fmt.Sprintf("%d", lineNo)
		}
		// Find the key and setup record to save
		for i, val := range row {
			if i < len(fieldNames) {
				fieldName = fieldNames[i]
				if idCol == i {
					key, err = ValueInterfaceToString(val)
					if err != nil {
						key = ""
					}
				}
			} else {
				fieldName = fmt.Sprintf(fmtColumnName, i+1)
			}
			record[fieldName] = val
		}
		if len(key) > 0 && len(record) > 0 {
			if c.HasKey(key) == true {
				if overwrite == true {
					err = c.Update(key, record)
					if err != nil {
						return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
					}
				} else if verboseLog == true {
					log.Printf("Skipped row %d, key %s exists in %s", lineNo, key, c.Name)
				}
			} else {
				err = c.Create(key, record)
				if err != nil {
					return lineNo, fmt.Errorf("can't write %+v to %s, %s", record, key, err)
				}
			}
		}
		if verboseLog == true && (lineNo%1000) == 0 {
			log.Printf("%d rows processed", lineNo)
		}
	}
	return lineNo, nil
}

func colToString(cell interface{}) string {
	var s string
	switch cell.(type) {
	case string:
		s = fmt.Sprintf("%s", cell)
	case json.Number:
		s = fmt.Sprintf("%s", cell.(json.Number).String())
	default:
		src, err := json.Marshal(cell)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		s = fmt.Sprintf("%s", src)
	}
	return s
}

// ExportCSV takes a reader and frame and iterates over the objects
// generating rows and exports then as a CSV file
func (c *Collection) ExportCSV(fp io.Writer, eout io.Writer, f *DataFrame, verboseLog bool) (int, error) {
	//, filterExpr string, dotExpr []string, colNames []string, verboseLog bool) (int, error) {
	keys := f.Keys[:]
	dotExpr := f.DotPaths
	colNames := f.Labels

	// write out colNames
	w := csv.NewWriter(fp)
	if err := w.Write(colNames); err != nil {
		return 0, err
	}

	var (
		cnt           int
		row           []string
		readErrors    int
		writeErrors   int
		dotpathErrors int
	)
	for i, key := range keys {
		data := map[string]interface{}{}
		if err := c.Read(key, data); err == nil {
			// write row out.
			row = []string{}
			for _, colPath := range dotExpr {
				col, err := dotpath.Eval(colPath, data)
				if err == nil {
					row = append(row, colToString(col))
				} else {
					if verboseLog == true {
						log.Printf("error in dotpath %q for key %q in %s, %s\n", colPath, key, c.workPath, err)
					}
					dotpathErrors++
					row = append(row, "")
				}
			}
			if err := w.Write(row); err == nil {
				cnt++
			} else {
				if verboseLog == true {
					log.Printf("error writing row %d from %s key %q, %s\n", i+1, c.workPath, key, err)
				}
				writeErrors++
			}
			data = nil
		} else {
			log.Printf("error reading %s %q, %s\n", c.workPath, key, err)
			readErrors++
		}
	}
	if readErrors > 0 || writeErrors > 0 || dotpathErrors > 0 && verboseLog == true {
		log.Printf("warning %d read error, %d write errors, %d dotpath errors in CSV export from %s", readErrors, writeErrors, dotpathErrors, c.workPath)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return cnt, err
	}
	return cnt, nil
}

// ExportTable takes a reader and frame and iterates over the objects
// generating rows and exports then as a CSV file
func (c *Collection) ExportTable(eout io.Writer, f *DataFrame, verboseLog bool) (int, [][]interface{}, error) {
	keys := f.Keys[:]
	dotExpr := f.DotPaths
	colNames := f.Labels

	var (
		cnt           int
		row           []interface{}
		readErrors    int
		dotpathErrors int
	)
	table := [][]interface{}{}
	// Copy column names to table
	for _, colName := range colNames {
		row = append(row, colName)
	}
	table = append(table, row)

	for _, key := range keys {
		data := map[string]interface{}{}
		if err := c.Read(key, data); err == nil {
			// write row out.
			row = []interface{}{}
			for _, colPath := range dotExpr {
				col, err := dotpath.Eval(colPath, data)
				if err == nil {
					row = append(row, col)
				} else {
					if verboseLog == true {
						log.Printf("error in dotpath %q for key %q in %s, %s\n", colPath, key, c.workPath, err)
					}
					dotpathErrors++
					row = append(row, nil)
				}
			}
			table = append(table, row)
			cnt++
			data = nil
		} else {
			log.Printf("error reading %s %q, %s\n", c.workPath, key, err)
			readErrors++
		}
	}
	if (readErrors > 0 || dotpathErrors > 0) && verboseLog == true {
		log.Printf("warning %d read error, %d dotpath errors in table export from %s", readErrors, dotpathErrors, c.workPath)
	}
	return cnt, table, nil
}

