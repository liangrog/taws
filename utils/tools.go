package utils

import (
	"fmt"
	"os"
)

// Check if providing key value exist in AWS resource tags
func HasTag(key, value string, tags []map[string]string) bool {
	for _, el := range tags {
		if el["Key"] == key &&
			el["Value"] == value {
			return true
		}
	}

	return false
}

// Exist with non-zero code
func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}
