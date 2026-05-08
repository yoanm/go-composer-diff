package compdiff_test

import (
	"encoding/json"
	"fmt"
	"testing"

	compdiff "github.com/yoanm/go-composer-diff"
)

func BenchmarkDiff_ComposerDiff(b *testing.B) {
	lockPrevious := generateComposerLockFile(1000)
	lockCurrent := generateComposerLockFile(1000)
	reqPrevious := generateComposerReqFile(1000)
	reqCurrent := generateComposerReqFile(1000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // pb.Next() returns false when the benchmark should stop
			_, err := compdiff.Diff(
				&compdiff.Input{
					Lock:        lockPrevious,
					Requirement: reqPrevious,
				},
				&compdiff.Input{
					Lock:        lockCurrent,
					Requirement: reqCurrent,
				},
			)
			if err != nil {
				b.Fatalf("Diff failed: %v", err)
			}
		}
	})
}

func generateComposerLockFile(count int) []byte {
	type pkg struct {
		Name      string            `json:"name"`
		Version   string            `json:"version"`
		Source    map[string]string `json:"source,omitempty"`
		Abandoned bool              `json:"abandoned"`
	}

	type lock struct {
		Packages []*pkg `json:"packages"`
	}

	lockObject := lock{
		Packages: make([]*pkg, count),
	}

	for i := range count {
		lockObject.Packages[i] = &pkg{
			Name:      fmt.Sprintf("vendor/package-%d", i),
			Version:   fmt.Sprintf("%d.%d.0", i%10, i%5),
			Source:    map[string]string{"reference": fmt.Sprintf("abc%d", i)},
			Abandoned: false,
		}
	}

	data, _ := json.Marshal(lockObject) //nolint:errchkjson // Not the purpose of the test

	return data
}

func generateComposerReqFile(count int) []byte {
	type reqStruct struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}

	reqObject := reqStruct{
		Require:    map[string]string{},
		RequireDev: map[string]string{},
	}

	for i := range count {
		if i%2 == 0 {
			reqObject.Require[fmt.Sprintf("vendor/package-%d", i)] = fmt.Sprintf("^%d.%d", i%10, i%5)
		} else {
			reqObject.RequireDev[fmt.Sprintf("vendor/package-%d", i)] = fmt.Sprintf("^%d.%d", i%10, i%5)
		}
	}

	data, _ := json.Marshal(reqObject) //nolint:errchkjson // Not the purpose of the test

	return data
}
