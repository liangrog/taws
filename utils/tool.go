package utils

import (
	"fmt"
	"os"
)

// Exist with non-zero code
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}
