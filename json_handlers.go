package dataset

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JSONUnmarshal is a custom JSON decoder so we can treat numbers easier
func JSONUnmarshal(src []byte, data interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	dec.UseNumber()
	err := dec.Decode(&data)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// JSONMarshal provides provide a custom json encoder to solve a an issue with
// HTML entities getting converted to UTF-8 code points by json.Marshal(), json.MarshalIndent().
func JSONMarshal(data interface{}) ([]byte, error) {
	buf1 := []byte{}
	w1 := bytes.NewBuffer(buf1)
	enc := json.NewEncoder(w1)
	enc.SetIndent("", "")
	enc.SetEscapeHTML(false)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	src1 := w1.Bytes()

	// compact the record so it takes up only one line.
	buf2 := []byte{}
	w2 := bytes.NewBuffer(buf2)
	err = json.Compact(w2, src1)
	src2 := w2.Bytes()
	return src2, err
}

// JSONMarshalIndent provides provide a custom json encoder to solve a an issue with
// HTML entities getting converted to UTF-8 code points by json.Marshal(), json.MarshalIndent().
func JSONMarshalIndent(data interface{}, prefix string, indent string) ([]byte, error) {
	buf := []byte{}
	w := bytes.NewBuffer(buf)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), err
}

// JSONIndent takes an byte slice of JSON source and returns an indented version.
func JSONIndent(src []byte, prefix string, indent string) []byte {
	buf := []byte{}
	w := bytes.NewBuffer(buf)
	json.Indent(w, src, prefix, indent)
	return w.Bytes()
}

// Dump takes an existing dataset collection and renders JSON objects
// one per line (i.e. JSONL, see https://jsonlines.org). The object structure
// written to the out buffer uses simple schema of a key attribute and a
// object attribute. This is regardles of the storage type of the collection
// being dumped.
//
// Here is an example of a single object being dump. The object key is
// "mirtle", the object is `{"one": 1}`.

// ```jsonl
//
//	{"key": "mirtle", "object": { "one": "1 }}
//
// ```
//
// Here is how you would use Dump in a Go project.
//
// ```go
//
//	cName := "mycollection.ds"
//	c, err := dataset.Open(cName)
//	if err != nil {
//	    ... // handle error
//	}
//	defer c.Close()
//	err := c.Dump(os.Stdout)
//	if err != nil {
//	    ... // handle error
//	}
//
// ```
func (c *Collection) Dump(out io.Writer) error {
	// Use a streaming approach if possible, e.g., iterate over the collection without loading all keys at once.
	// If not possible, proceed with the current approach but ensure minimal memory usage.
	keys, err := c.Keys()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return fmt.Errorf("collection %q is empty", c.Name)
	}

	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "") // For compact JSONL output

	errCnt := 0
	tot := len(keys)
	for i, key := range keys {
		obj := map[string]interface{}{}
		err := c.Read(key, obj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING (%d/%d) failed to read %q from %q, %s\n", i, tot, key, c.Name, err)
			errCnt++
			continue
		}

		rec := map[string]interface{}{
			"key":    key,
			"object": obj,
		}

		if err := enc.Encode(rec); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING (%d/%d) failed to encode %q from %q, %s\n", i, tot, key, c.Name, err)
			errCnt++
			continue
		}
	}

	if errCnt > 0 {
		return fmt.Errorf("%d dump errors for %q", errCnt, c.Name)
	}
	return nil
}

// Load reads JSONL from an io.Reader. The JSONL object should have two attributes.
// The first is "key" should should be a unique string the object is "object" which
// is the JSON object to be stored in the collection. The collection needs to exist.
// If the overwrite parameter is set to true then the object read will overwrite
// any objects with the same key. If overwrite is false you will get a warning mesage
// that the object was skipped due to duplicate key.
// The third parameter is the size of the input buffer scanned in megabytes. If
//
//	the value is less or equal to zero then it defaults to 1 megabyte buffer.
//
// ```
//
//		 cName := "mycollection.ds"
//		 c, err := dataset.open(cName)
//		 if err != nil {
//		    // ... handle error
//		 }
//		 defer c.Close()
//	  // use the default buffer size
//		 err = c.Load(os.Stdin, maxCapacity, 0)
//		 if err != nil {
//		    // ... handle error
//		 }
//
// ```
func (c *Collection) Load(in io.Reader, overwrite bool, maxCapacity int) error {
	// Set a large buffer size if maxCapacity is specified
	bufSize := 1024 * 1024 // Default: 1 MB
	if maxCapacity > 0 {
		bufSize = maxCapacity * 1024 * 1024
	}

	scanner := bufio.NewScanner(in)
	// Use a custom split function to handle very large lines
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		// Look for a newline character
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, return the remaining data
		if atEOF {
			return len(data), data, nil
		}
		// Request more data
		return 0, nil, nil
	})
	// Set the buffer size
	buf := make([]byte, bufSize)
	scanner.Buffer(buf, bufSize)

	dec := json.NewDecoder(nil)
	errCnt := 0
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		dec = json.NewDecoder(bytes.NewReader(line))
		dec.UseNumber()

		var rec map[string]interface{}
		if err := dec.Decode(&rec); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to decode line %d for %q, %s\n", lineNum, c.Name, err)
			errCnt++
			continue
		}

		key, keyOk := rec["key"].(string)
		obj, objOk := rec["object"].(map[string]interface{})

		if !keyOk || !objOk {
			fmt.Fprintf(os.Stderr, "WARNING (line %d): missing key or object\n", lineNum)
			errCnt++
			continue
		}

		keyExists := c.HasKey(key)
		if keyExists {
			if overwrite {
				if err := c.Update(key, obj); err != nil {
					fmt.Fprintf(os.Stderr, "WARNING (line %d): failed to update %q, %s\n", lineNum, key, err)
					errCnt++
				}
			} else {
				fmt.Fprintf(os.Stderr, "WARNING (line %d): duplicate key %q, skipping\n", lineNum, key)
				errCnt++
			}
		} else {
			if err := c.Create(key, obj); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING (line %d): failed to create %q, %s\n", lineNum, key, err)
				errCnt++
			}
		}
	}

	// Handle scanner errors (e.g., "token too long")
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: scanning error at line %d, %s\n", lineNum, err)
		errCnt++
	}

	if errCnt > 0 {
		return fmt.Errorf("%d load errors for %q", errCnt, c.Name)
	}
	return nil
}
