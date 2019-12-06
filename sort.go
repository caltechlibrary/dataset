//
// Package dataset includes the operations needed for processing collections of JSON documents and their attachments.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu> and Tom Morrel, <tmorrell@library.caltech.edu>
//
// Copyright (c) 2019, Caltech
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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/dotpath"
)

//KeyValue holds an ID string and value interface, this lets us work with numeric keys and to sort them.
type KeyValue struct {
	// JSON Record ID in collection
	ID string
	// The value of the field to be sorted from record
	Value interface{}
}

//KeyValues is a list of keys (strings) to records. This type exists to allow easy sorting.
type KeyValues []KeyValue

func (a KeyValues) Len() int {
	return len(a)
}

func (a KeyValues) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a KeyValues) Less(i, j int) bool {
	switch a[i].Value.(type) {
	case string:
		return a[i].Value.(string) < a[j].Value.(string)
	case float64:
		return a[i].Value.(float64) < a[j].Value.(float64)
	case int64:
		return a[i].Value.(int64) < a[j].Value.(int64)
	case int:
		return a[i].Value.(int) < a[j].Value.(int)
	case json.Number:
		return a[i].Value.(json.Number) < a[j].Value.(json.Number)
	default:
		s1 := fmt.Sprintf("%+v", a[i].Value)
		s2 := fmt.Sprintf("%+v", a[j].Value)
		return s1 < s2
	}
}

// sortBy takes a listof record ids in a collection, performs a single field sort for a given
// dotpath (dpath) returning a sorted list of ids.
func (c *Collection) sortBy(ids []string, dpath string, ascending bool) ([]string, error) {
	var (
		t1 string
		t2 string
	)
	if len(ids) == 0 || len(dpath) == 0 {
		return ids, nil
	}
	// FIXME: loop through ids, collect the values from based on dotpath
	// sort the values
	list := []KeyValue{}
	for _, id := range ids {
		// Get the record for given id
		rec := map[string]interface{}{}
		if err := c.Read(id, rec, false); err == nil {
			// if exist extract the value
			if value, err := dotpath.Eval(dpath, rec); err == nil {
				// Verify we have same type on data or return error
				sv := KeyValue{
					ID:    id,
					Value: value,
				}
				if len(list) == 0 {
					t1 = fmt.Sprintf("%T", sv.Value)
				} else {
					t2 = fmt.Sprintf("%T", sv.Value)
					if t1 != t2 {
						return ids, fmt.Errorf("value for %s (%s) does not match type %s", id, t2, t1)
					}
				} // add it to list of SortableValue
				list = append(list, sv)
			}
		}
	}
	// Sort the values in list of SortableValues
	sort.Sort(KeyValues(list))
	// Return the id list with sorted values in either ascending or decending order
	ids = []string{}
	if ascending == true {
		for _, sv := range list {
			ids = append(ids, sv.ID)
		}
	} else {
		for i := len(list) - 1; i >= 0; i-- {
			ids = append(ids, list[i].ID)
		}
	}
	return ids, nil
}

// KeySortByExpression takes a array of keys and a sort expression and turns a sorted list of keys.
func (c *Collection) KeySortByExpression(keys []string, expr string) ([]string, error) {
	ascending := true
	if strings.HasPrefix(expr, "-") {
		ascending = false
		expr = expr[1:]
	}
	if strings.HasPrefix(expr, "+") {
		expr = expr[1:]
	}
	return c.sortBy(keys, expr, ascending)
}
