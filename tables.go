//
// table.go provides some utility functions to move string one and
// two dimensional slices into/out of one and two dimensional slices.
//
package dataset

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

// ValueInterfaceToString - takes a interface{} and renders it as a string
func ValueInterfaceToString(val interface{}) (string, error) {
	switch val.(type) {
	case string:
		return val.(string), nil
	case json.Number:
		return val.(json.Number).String(), nil
	case int:
		return fmt.Sprintf("%d", val), nil
	case int64:
		return fmt.Sprintf("%d", val), nil
	case float64:
		return fmt.Sprintf("%f", val), nil
	case rune:
		return fmt.Sprintf("%s", val), nil
	case byte:
		return fmt.Sprintf("%d", val), nil
	case []byte:
		return fmt.Sprintf("%s", val), nil
	case []rune:
		return fmt.Sprintf("%s", val), nil
	default:
		src, err := JSONMarshal(val)
		if err != nil {
			return "", fmt.Errorf("unknown type conversion, %T, %s", val, err)
		}
		return fmt.Sprintf("%s", src), nil
	}
}

var (
	reInteger = regexp.MustCompile(`[0-9]+`)
	reReal    = regexp.MustCompile(`[0-9]+\.[0-9]+`)
	reBool    = regexp.MustCompile(`[tT][rR][uU][eE]|[fF][aA][lL][sS][eE]|0|1`)
)

// ValueStringToInterface takes a string and returns an interface{}
func ValueStringToInterface(s string) (interface{}, error) {
	// FIXME: Need to introduce some sort of smart data conversion
	// so that a string representation of a number becomes a json.Number
	// I need to try to derive type and then
	// create a variable of that type and add it to the
	// interface{}
	switch {
	case reInteger.MatchString(s):
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return interface{}(i), nil
	case reReal.MatchString(s):
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return interface{}(f), nil
	case reBool.MatchString(s):
		b, err := strconv.ParseBool(s)
		if err != nil {
			return nil, err
		}
		return interface{}(b), nil
	}
	return interface{}(s), nil
}

// RowStringToInterface takes a 1D slice of string and returns
// a 1D slice of interface{}
func RowStringToInterface(r []string) []interface{} {
	row := []interface{}{}
	for _, cell := range r {
		val, err := ValueStringToInterface(cell)
		if err == nil {
			row = append(row, val)
		} else {
			row = append(row, cell)
		}
	}
	return row
}

// RowInterfaceToString takes a 1D slice of interface{} and
// returns a 1D slice of string, of conversion then cell
// will be set to empty string.
func RowInterfaceToString(r []interface{}) []string {
	cells := []string{}
	for _, cell := range r {
		val, err := ValueInterfaceToString(cell)
		if err == nil {
			cells = append(cells, val)
		} else {
			cells = append(cells, "")
		}
	}
	return cells
}

// TableStringToInterface takes a 2D slice of string and returns
// an 2D slice of interface{}.
func TableStringToInterface(t [][]string) [][]interface{} {
	table := [][]interface{}{}
	for _, row := range t {
		table = append(table, RowStringToInterface(row))
	}
	return table
}

// TableInterfaceToString takes a 2D slice of interface{}
// holding simple types (e.g. string, int, int64, float, float64,
// rune) and returns a 2D slice of string suitable for working
// with the csv encoder package.  Uses ValueInterfaceToString()
// for conversion storing an empty string if they is an error.
func TableInterfaceToString(t [][]interface{}) [][]string {
	table := [][]string{}
	for _, row := range t {
		table = append(table, RowInterfaceToString(row))
	}
	return table
}

