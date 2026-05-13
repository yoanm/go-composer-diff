package compdiff_test

import (
	"fmt"
	"os"
	"testing"

	compdiff "github.com/yoanm/go-composer-diff"
)

func Benchmark_BuildMapFromBytes(b *testing.B) {
	// Load fixture files
	reqData, err := os.ReadFile("./testdata/complex-composer.json")
	if err != nil {
		b.Fatal(fmt.Errorf("error while reading requirement file = %w", err))
	}

	lockData, err := os.ReadFile("./testdata/complex-composer.lock")
	if err != nil {
		b.Fatal(fmt.Errorf("error while reading lock file = %w", err))
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // pb.Next() returns false when the benchmark should stop
			if _, err2 := compdiff.BuildMapFromBytes(reqData, lockData); err2 != nil {
				b.Fatalf("BuildMapFromBytes failed: %v", err2)
			}
		}
	})
}
