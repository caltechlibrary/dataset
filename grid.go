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
	"log"
	"os"

	// Caltech Library Packages
	"github.com/caltechlibrary/dotpath"
)

// Grid takes a set of collection keys and builds a grid (a 2D array cells)
// from the array of keys and dot paths provided
func (c *Collection) Grid(keys []string, dotPaths []string, verbose bool) ([][]interface{}, error) {
	pid := os.Getpid()
	rows := make([][]interface{}, len(keys))
	col_cnt := len(dotPaths)
	for i, key := range keys {
		rec := map[string]interface{}{}
		err := c.Read(key, rec)
		if err != nil {
			return nil, err
		}
		rows[i] = make([]interface{}, col_cnt)
		for j, dpath := range dotPaths {
			value, err := dotpath.Eval(dpath, rec)
			if err == nil {
				rows[i][j] = value
			} else if verbose == true {
				log.Printf("(pid: %d) WARNING: skipped %s for cell %d row %d, %s", pid, dpath, j, i, err)
			}
		}
		if verbose && (i > 0) && ((i % 1000) == 0) {
			log.Printf("(pid: %d) %d keys processed", pid, i)
		}
	}
	return rows, nil
}
