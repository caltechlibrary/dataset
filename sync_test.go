package dataset

import (
	"bytes"
	"encoding/csv"
	"testing"
)

func TestTableColumnBehavior(t *testing.T) {
	src := []byte(`
first, second,third, fourth
1,,3,
2,2,2,
3,3,
`)
	r := csv.NewReader(bytes.NewBuffer(src))
	r.FieldsPerRecord = -1
	csvTable, err := r.ReadAll()
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	table := TableStringToInterface(csvTable)

	expectedObjs := map[string]map[string]interface{}{}
	expectedObjs["1"] = map[string]interface{}{
		"key":   "1",
		"first": 1,
		"third": 3,
	}
	expectedObjs["2"] = map[string]interface{}{
		"key":    "2",
		"first":  2,
		"second": 2,
		"third":  2,
	}
	expectedObjs["3"] = map[string]interface{}{
		"key":    "3",
		"first":  3,
		"second": 3,
	}

	for i, row := range table {
		if i > 0 {
			key, err := ValueInterfaceToString(row[0])
			if err != nil {
				t.Errorf("expected (%T) %+v to be convertable to string, %s", row[0], row[0], err)
				continue
			}
			for j, fieldName := range []string{"first", "second", "third", "fourth"} {
				if eObj, ok := expectedObjs[key]; ok == true {
					obj := eObj
					if len(obj) > len(row) {
						t.Errorf("row %d is short cells, obj %+v, row %+v", i, obj, row)
					}
					if val, ok := obj[fieldName]; ok == true && val != row[j] {
						t.Errorf("row %d, col %d, expected (%T) %+v, got (%T) %+v", i, j, val, val, row[j], row[j])
					}
				} else {
					t.Errorf("unexpected key %q in row %d table %+v\n", key, i, table)
				}
			}
		}
	}
}

