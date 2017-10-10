package cmd

import (
	"strings"
)

// Allowing the same argument multiple time
type MultiStrVar []string

// Presenting a comma-separated string
func (v *MultiStrVar) String() string {
	return strings.Join(*v, ",")
}

// Append multiple string vars to the slice
func (v *MultiStrVar) Set(s string) error {
	if *v == nil {
		*v = make([]string, 0, 1)
	}

	*v = append(*v, s)
	return nil
}
