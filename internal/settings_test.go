package internal_test

import (
	"aproc/internal"
	"testing"
)

func TestCreateEmptySettings(t *testing.T) {
	t.Parallel()

	if err := internal.CreateEmptySettings(); err != nil {
		t.Fatal(err)
	}
}
