package lib_test

import (
	"aproc/lib"
	"testing"
)

func TestCreateEmptySettings(t *testing.T) {
	t.Parallel()

	if err := lib.CreateEmptySettings(); err != nil {
		t.Fatal(err)
	}
}
