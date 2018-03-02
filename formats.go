//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package dataset

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"text/template"

	// Caltech Library packages
	"github.com/caltechlibrary/mkpage"

	// 3rd Party packages
	"github.com/blevesearch/bleve"
)

// JSONFormatter writes out JSON representation using encoding/json
func JSONFormatter(out io.Writer, results *bleve.SearchResult, prettyPrint bool) error {
	var (
		src []byte
		err error
	)
	if prettyPrint {
		src, err = json.MarshalIndent(results, "", "    ")
	} else {
		src, err = json.Marshal(results)
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s\n", src)
	return nil
}

// CSVFormatter writes out CSV representation using encoding/csv
func CSVFormatter(out io.Writer, results *bleve.SearchResult, colNames []string, skipHeaderRow bool) error {
	// Note: we need to provide the fieldnames that will be come columns
	w := csv.NewWriter(out)
	// write a header row if needed
	if skipHeaderRow == false {
		if err := w.Write(colNames); err != nil {
			return err
		}
	}
	for _, hit := range results.Hits {
		row := []string{}
		for _, col := range colNames {
			if col == "_id" {
				row = append(row, hit.ID)
			} else if col == "_index" {
				row = append(row, path.Base(hit.Index))
			} else if val, ok := hit.Fields[col]; ok == true {
				switch val := val.(type) {
				case int:
					row = append(row, strconv.FormatInt(int64(val), 10))
				case uint:
					row = append(row, strconv.FormatUint(uint64(val), 10))
				case int64:
					row = append(row, strconv.FormatInt(val, 10))
				case uint64:
					row = append(row, strconv.FormatUint(val, 10))
				case float64:
					row = append(row, strconv.FormatFloat(val, 'G', -1, 64))
				case string:
					row = append(row, strings.TrimSpace(val))
				default:
					row = append(row, strings.TrimSpace(fmt.Sprintf("%s", val)))
				}
			} else {
				row = append(row, "")
			}
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

// Formatter writes out a format based on the specified template name merging any additional pageData provided
func Formatter(out io.Writer, results *bleve.SearchResult, tmpl *template.Template, tName string, pageData map[string]string) error {
	src, err := json.Marshal(results)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	decoder := json.NewDecoder(bytes.NewBuffer(src))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	if len(pageData) > 0 {
		if pData, err := mkpage.ResolveData(pageData); err != nil {
			return fmt.Errorf("Can't resolve data source %s", err)
		} else {
			// Merge the k/v for the page if not supplied by results
			for k, v := range pData {
				if _, ok := data[k]; ok == false {
					data[k] = v
				}
			}
		}
	}
	return tmpl.ExecuteTemplate(out, tName, data)
}
