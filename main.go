package main

import (
	"flag"
	"fmt"
	"slices"
	"strings"
)

var (
	target          = flag.Int64("target", 0, "-target <used memory (GB)>")
	n               = flag.Int("n", 20, "[-n <no of results. defaults to 20>]")
	excludedKeysStr = flag.String("e", "Pages free", `-e "Pages free"`)
	includedKeysStr = flag.String("i", "Pages active,Pages wired down", `-i "Pages active,Pages wired down"`)
)

func main() {
	flag.Parse()

	fmt.Printf("Starting with target %d...\n", *target)

	targetBytes := *target << 20

	mi, err := meminfo()
	if err != nil {
		panic(err)
	}

	combinations := calculateCombination(&mi, targetBytes)

	slices.SortFunc(combinations, func(c1, c2 Combination) int {
		return int(c1.diff - c2.diff)
	})

	for i, comb := range combinations {
		if i < *n {
			fmt.Printf("%d MB diff, %#v, %d MB total\n", comb.diff>>20, comb.keys, comb.byteSize>>20)
		}
	}
}

func calculateCombination(mi *MemInfo, targetBytes int64) []Combination {
	keys := make([]string, 0, len(mi.values))
	excludedKeys := strings.Split(*excludedKeysStr, ",")
	includedKeys := strings.Split(*includedKeysStr, ",")

	for k, v := range mi.values {
		if !slices.Contains(excludedKeys, k) && !slices.Contains(includedKeys, k) && v > 0 {
			keys = append(keys, k)
		}
	}

	var currentSize int64
	for _, k := range includedKeys {
		currentSize += mi.values[k] * mi.pageSize
	}
	combinations := []Combination{
		{
			byteSize: currentSize,
			diff:     abs(targetBytes - currentSize),
			keys:     includedKeys,
		},
	}

	var f func(int64, int, []string)
	f = func(current int64, index int, currentCombination []string) {
		v := mi.values[keys[index]]

		if v == 0 {
			return
		}

		size := v * mi.pageSize

		currentCombination = append(currentCombination, keys[index])

		combinations = append(combinations, Combination{
			keys:     currentCombination,
			byteSize: current + size,
			diff:     abs(targetBytes - (current + size)),
		})

		for i := index + 1; i < len(keys); i += 1 {
			f(current+size, i, slices.Clone(currentCombination))
		}
	}

	for i := range keys {
		f(currentSize, i, slices.Clone(includedKeys))
	}

	return combinations
}

func abs(i int64) int64 {
	if i < 0 {
		return i * -1
	}
	return i
}

type Combination struct {
	keys     []string
	byteSize int64
	diff     int64
}
