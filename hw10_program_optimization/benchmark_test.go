package hw10programoptimization

import (
	"archive/zip"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func BenchmarkCaseFromFrozenTest(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// much better to seek start, but zip file does't support this(
		r, err := zip.OpenReader("testdata/users.dat.zip")
		require.NoError(b, err)
		data, err := r.File[0].Open()
		require.NoError(b, err)

		b.StartTimer()
		_, err = GetDomainStat(data, "biz")
		b.StopTimer()
		require.NoError(b, err)

		r.Close()
	}
}

func BenchmarkRandomData(b *testing.B) {
	b.ReportAllocs()
	rdr := strings.NewReader(generateRandomInput(1000))
	suffix := gofakeit.DomainSuffix()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetDomainStat(rdr, suffix)
	}
}

func generateRandomInput(limit int) string {
	gofakeit.Seed(time.Now().UnixMicro())
	result := strings.Builder{}
	for i := 0; i < limit; i++ {
		result.WriteString(fmt.Sprintf(`{"Name":"%s", "Email":"%s"}\n`, gofakeit.Name(), gofakeit.Email()))
	}
	line := result.String()

	return line
}
