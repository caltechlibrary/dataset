package dataset

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type DSImport struct {
	Comma string
	Comment string
	Overwrite bool
	LazyQuotes bool
	TrimLeadingSpace bool
}

// normalizeDelimiters handles the messy translation from a format string
// received as an option in the cli to something useful to pass to Join.
func normalizeDelimiter(s string) string {
        if strings.Contains(s, `\n`) {
                s = strings.Replace(s, `\n`, "\n", -1)
        }
        if strings.Contains(s, `\t`) {
                s = strings.Replace(s, `\t`, "\t", -1)
        }
        return s
}

// normalizeDelimiterRune take a delimiter string and returns a single Rune
func normalizeDelimiterRune(s string) rune {
        runes := []rune(normalizeDelimiter(s))
        if len(runes) > 0 {
                return runes[0]
        }
        return ','
}


func (app *DSImport) Run(in io.Reader, out io.Writer, eout io.Writer, cName string, keyColumn string) error {
	c, err := Open(cName)
	if err != nil {
		return err
	}
	defer c.Close()

	r := csv.NewReader(in)
	if app.Comma != "" {
		r.Comma = normalizeDelimiterRune(app.Comma)
	}
	if app.Comment != "" {
		r.Comment = normalizeDelimiterRune(app.Comment)
	}
	r.LazyQuotes = app.LazyQuotes
    r.TrimLeadingSpace = app.TrimLeadingSpace

	i := 0
	header := []string{}
	keyIndex := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(eout, "row %d: %s\n", i, err)
			continue
		}
		if i == 0 {
			header = row
			for i, val := range header {
				if strings.Compare(keyColumn, val) == 0 {
					keyIndex = i
					break
				}
			}
		} else {
			if keyIndex < len(row) {
				key := row[keyIndex]
				values := map[string]interface{}{}
				for j, attr := range header {
					if j < len(row) {
						values[attr] = row[j]
					} else {
						fmt.Fprintf(eout, "row %d: can't find column (%d) %q\n", j, attr)
					}
				}
				if c.HasKey(key) {
					if app.Overwrite {
						if err := c.Update(key, values); err != nil  {
							fmt.Fprintf(eout, "row %d: failed to update %q in collection, %s\n", i, key, err)
						}
					} else {
						fmt.Fprintf(eout, "row %d: key already exists in collection %q\n", i, key)
					}
				} else {
					if err := c.Create(key, values); err != nil  {
						fmt.Fprintf(eout, "row %d: failed to create %q in collection, %s\n", i, key, err)
					}
				}
			} else {
				fmt.Fprint(eout, "row %d: can't find key column", i)
			}
		}
		i++
	}

	return nil
}
