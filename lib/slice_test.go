package lib_test

import (
	"aproc/lib"
	"testing"
)

func TestSlice(t *testing.T) {
	t.Parallel()

	slice1 := []string{"1", "2", "3", "6", "8"}
	slice2 := []string{"2", "3", "5", "0"}
	un := lib.Union(slice1, slice2)

	if !lib.StringSliceEqualBCE(un, []string{"1", "2", "3", "6", "8", "5", "0"}) {
		t.Fatalf("slice1 并 slice2 的结果为 %v", un)
	}

	in := lib.Intersect(slice1, slice2)

	if !lib.StringSliceEqualBCE(in, []string{"2", "3"}) {
		t.Fatalf("slice1 交 slice2 的结果为 %v", in)
	}

	di := lib.Difference(slice1, slice2)

	if !lib.StringSliceEqualBCE(di, []string{"1", "6", "8"}) {
		t.Fatalf("slice1 - slice2 的结果为 %v", di)
	}
}
