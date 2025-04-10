package hw10programoptimization

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func generateTestData(count int) []byte {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString("{\"Id\":")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString(",\"Email\":\"user")
		sb.WriteString(strconv.Itoa(i))
		sb.WriteString("@example.com\"}\n")
	}
	return []byte(strings.TrimSuffix(sb.String(), "\n"))
}

func BenchmarkGetDomainStatOld(b *testing.B) {
	testData := generateTestData(1000)
	domain := "example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(testData)
		_, err := GetDomainStatOld(r, domain)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetDomainStat(b *testing.B) {
	testData := generateTestData(1000)
	domain := "example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(testData)
		_, err := GetDomainStat(r, domain)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// $ go test -bench=. -benchmem -benchtime=3s -v
// === RUN   TestGetDomainStat
// === RUN   TestGetDomainStat/find_'com'
// === RUN   TestGetDomainStat/find_'gov'
// === RUN   TestGetDomainStat/find_'unknown'
// --- PASS: TestGetDomainStat (0.00s)
//     --- PASS: TestGetDomainStat/find_'com' (0.00s)
//     --- PASS: TestGetDomainStat/find_'gov' (0.00s)
//     --- PASS: TestGetDomainStat/find_'unknown' (0.00s)
// goos: linux
// goarch: amd64
// pkg: github.com/BukhalovaK/otuDZ/hw10_program_optimization
// cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz
// BenchmarkGetDomainStatOld
// BenchmarkGetDomainStatOld-8           16         196577254 ns/op        243014404 B/op   2506025 allocs/op
// BenchmarkGetDomainStat
// BenchmarkGetDomainStat-8           16728            218487 ns/op          157720 B/op       2014 allocs/op
// PASS
// ok      github.com/BukhalovaK/otuDZ/hw10_program_optimization   9.196s
