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

	for i := 0; i < b.N; i++ {
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		data, _ := r.File[0].Open()

		b.StartTimer()
		GetDomainStat(data, "biz")
		b.StopTimer()

		data.Close()
	}
}

func BenchmarkRandomData(b *testing.B) {
	b.StopTimer()
	lines := generateRandomBuffer(1000)
	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GetDomainStat(bytes.NewBufferString(lines), gofakeit.DomainSuffix())
	}
	b.StopTimer()
}

func BenchmarkExtractEmail(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		gofakeit.Seed(time.Now().UnixMicro())
		info := gofakeit.Person()

		if bytes, e := jsoniter.Marshal(*info); e == nil {
			b.StartTimer()
			extractEmail(bytes)
			b.StopTimer()
		} else {
			b.FailNow()
		}
	}
}

func BenchmarkHandleEmail(b *testing.B) {
	b.StopTimer()
	stat := DomainStat{}
	for i := 0; i < b.N; i++ {
		gofakeit.Seed(time.Now().UnixMicro())
		email := gofakeit.Email()

		b.StartTimer()
		handleEmail(email, "com", &stat)
		b.StopTimer()
	}
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
