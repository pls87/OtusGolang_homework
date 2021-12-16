package hw10programoptimization

import (
	"archive/zip"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func BenchmarkCaseFromFrozenTest(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		data, _ := r.File[0].Open()
		b.StartTimer()
		GetDomainStat(data, "biz")
		r.Close()
	}
}

func BenchmarkRandomData(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		rdr := strings.NewReader(generateRandomBuffer(1000))
		suffix := gofakeit.DomainSuffix()
		b.StartTimer()
		GetDomainStat(rdr, suffix)
	}
}

func generateRandomBuffer(limit int) string {
	gofakeit.Seed(time.Now().UnixMicro())
	result := strings.Builder{}
	for i := 0; i < limit; i++ {
		result.WriteString(fmt.Sprintf(`{"Name":"%s", "Email":"%s"}\n`, gofakeit.Name(), gofakeit.Email()))
	}
	line := result.String()

	return line
}
