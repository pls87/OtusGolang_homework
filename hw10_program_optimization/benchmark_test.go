package hw10programoptimization

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	jsoniter "github.com/json-iterator/go"
)

func BenchmarkCaseFromFrozenTest(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	data, _ := r.File[0].Open()

	b.StartTimer()
	GetDomainStat(data, "biz")
	b.StopTimer()

	data.Close()
}

func BenchmarkRandomData(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	bts := bytes.NewBufferString(generateRandomBuffer(1000))
	suffix := gofakeit.DomainSuffix()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		GetDomainStat(bts, suffix)
	}
	b.StopTimer()
}

func generateRandomBuffer(limit int) string {
	gofakeit.Seed(time.Now().UnixMicro())
	result := strings.Builder{}
	for i := 0; i < limit; i++ {
		if bytes, e := jsoniter.Marshal(*gofakeit.Person()); e == nil {
			result.Write(bytes)
			result.WriteString("\n")
		}
	}

	return result.String()
}
