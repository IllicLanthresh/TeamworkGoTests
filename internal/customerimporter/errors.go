package customerimporter

import (
	"fmt"
	"github.com/IllicLanthresh/TeamworkGoTests/pkg/helperTypes"
)

type KeyNotFoundError struct {
	key   string
	slice helperTypes.StringSlice
}

func (e KeyNotFoundError) Error() string {
	return fmt.Sprintf("couldn't find \"%s\" in %v", e.key, e.slice)
}

type CsvPathInvalidError struct {
	path string
}

func (e CsvPathInvalidError) Error() string {
	return fmt.Sprintf("\"%s\" is not a valid CSV file path", e.path)
}

type MissingEmailKey struct{}

func (_ MissingEmailKey) Error() string {
	return "you need to specify an email key"
}
