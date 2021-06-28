package customerimporter

import (
	"regexp"
	"testing"
)

func Test_emailRegex_compiles(t *testing.T) {
	_, err := regexp.Compile(emailRegex)
	if err != nil {
		t.Error(err)
	}
}
