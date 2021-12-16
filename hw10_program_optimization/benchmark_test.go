package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkCaseFromFrozenTest(b *testing.B) {
	r, _ := zip.OpenReader("testdata/users.dat.zip")
	data, _ := r.File[0].Open()

	b.StartTimer()
	GetDomainStat(data, "biz")
	b.StopTimer()
}
