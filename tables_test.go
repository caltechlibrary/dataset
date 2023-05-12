package dataset

import (
	"fmt"
	"strings"
	"testing"
)

func TestCellConversion(t *testing.T) {
	var (
		val interface{}
		err error
	)

	expectedS := "Hello World"
	val, err = ValueInterfaceToString([]byte(expectedS))
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if strings.Compare(expectedS, fmt.Sprintf("%s", val)) != 0 {
		t.Errorf("expected %q, got (%T) %+v\n", expectedS, val.(string), val.(string))
	}

	val, err = ValueStringToInterface(expectedS)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if strings.Compare(expectedS, fmt.Sprintf("%s", val)) != 0 {
		t.Errorf("expected %q, got (%T) %+v\n", expectedS, val.(string), val.(string))
	}

	expectedS = "1"
	expectedI := 1
	val, err = ValueStringToInterface(expectedS)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	if expectedI != val.(int) {
		t.Errorf("expected %d, got (%T) %+v\n", expectedI, val, val)
	}

}

