package cli

import (
	"fmt"
	"os"
)

// OnError writes an error message to out if err != nil
// taking into consideration the state of quiet
func OnError(out *os.File, err error, quiet bool) {
	if err != nil && quiet == false {
		fmt.Fprintf(out, "%s\n", err)
	}
}

// ExitOnError is used by the cli programs to
// handle exit cuasing errors constitantly.
// E.g. it respects the -quiet flag past to it.
func ExitOnError(out *os.File, err error, quiet bool) {
	if err != nil {
		if quiet == false {
			fmt.Fprintf(out, "%s\n", err)
		}
		os.Exit(1)
	}
}
