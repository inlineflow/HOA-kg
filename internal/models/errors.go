package models

import "fmt"

type PathValueParseError struct {
	ResourceKey string
	ParseError  error
}

func (e *PathValueParseError) Error() string {
	return fmt.Sprintf("Failed to parse %s from path value: %v", e.ResourceKey, e.ParseError)
}
