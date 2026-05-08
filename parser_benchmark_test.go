package compdiff_test

import (
	"encoding/json"
	"fmt"
	compdiff "github.com/yoanm/go-composer-diff"
	"testing"
)

func BenchmarkParseLock(b *testing.B) {
	data := generateLockFile(1000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // pb.Next() returns false when the benchmark should stop
			if _, err := compdiff.ParseLock(data); err != nil {
				b.Fatalf("ParseLock failed: %v", err)
			}
		}
	})
}

func generateLockFile(count int) []byte {
	type pkgStruct struct {
		Name      string            `json:"name"`
		Version   string            `json:"version"`
		Source    map[string]string `json:"source,omitempty"`
		Abandoned bool              `json:"abandoned"`
	}

	type lock struct {
		Packages    []*pkgStruct `json:"packages"`
		PackagesDev []*pkgStruct `json:"packages-dev"`
	}

	lockObject := lock{
		Packages:    []*pkgStruct{},
		PackagesDev: []*pkgStruct{},
	}

	for cnt := range count {
		pkg := pkgStruct{
			Name:      fmt.Sprintf("vendor/package-%d", cnt),
			Version:   fmt.Sprintf("%d.%d.0", cnt%10, cnt%5),
			Source:    map[string]string{"reference": fmt.Sprintf("abc%d", cnt)},
			Abandoned: false,
		}
		if cnt%2 == 0 {
			lockObject.Packages = append(lockObject.Packages, &pkg)
		} else {
			lockObject.PackagesDev = append(lockObject.PackagesDev, &pkg)
		}
	}

	data, _ := json.Marshal(lockObject) //nolint:errchkjson // Not the purpose of the test

	return data
}

func BenchmarkParseReq(b *testing.B) {
	data := generateReqFile(1000)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // pb.Next() returns false when the benchmark should stop
			if _, err := compdiff.ParseReq(data); err != nil {
				b.Fatalf("ParseReq failed: %v", err)
			}
		}
	})
}

func generateReqFile(count int) []byte {
	type reqStruct struct {
		Require    map[string]string `json:"require"`
		RequireDev map[string]string `json:"require-dev"`
	}

	reqObject := reqStruct{
		Require:    map[string]string{},
		RequireDev: map[string]string{},
	}

	for cnt := range count {
		key, value := fmt.Sprintf("vendor/package-%d", cnt), fmt.Sprintf("^%d.%d", cnt%10, cnt%5)
		if cnt%2 == 0 {
			reqObject.Require[key] = value
		} else {
			reqObject.RequireDev[key] = value
		}
	}

	data, _ := json.Marshal(reqObject) //nolint:errchkjson // Not the purpose of the test

	return data
}
