package main

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
)

func writeTempReport(content string) string {
	f, err := os.CreateTemp(GinkgoT().TempDir(), "report-*.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(content)
	return f.Name()
}
