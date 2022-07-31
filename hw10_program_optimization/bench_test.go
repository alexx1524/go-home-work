package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkGetDomainStat(b *testing.B) {
	readCloser, _ := zip.OpenReader("testdata/users.dat.zip")
	defer readCloser.Close()

	content, _ := readCloser.File[0].Open()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(content, "org")
	}
}
