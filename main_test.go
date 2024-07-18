package main

import (
	"slices"
	"testing"
)

func TestMain(t *testing.T) {
	mi := MemInfo{
		pageSize: 10,
		values: map[string]int64{
			"a": 1,
			"b": 2,
			"c": 3,
			"d": 4,
		},
	}

	combinations := calculateCombination(&mi, 90, []string{}, []string{})

	slices.SortFunc(combinations, func(c1, c2 Combination) int {
		return int(c1.diff - c2.diff)
	})

	if combinations[0].diff != 0 {
		t.Errorf("expected diff to be %d", 0)
	}
	if combinations[0].byteSize != 90 {
		t.Errorf("expected byteSize to be %d", 90)
	}
	if len(combinations[0].keys) != 3 {
		t.Errorf("expected length of keys to be %d", 3)
	} else if !slices.Contains(combinations[0].keys, "b") || !slices.Contains(combinations[0].keys, "c") || !slices.Contains(combinations[0].keys, "d") {
		t.Errorf("expected keys to be %v", []string{"b", "c", "d"})
	}
}
