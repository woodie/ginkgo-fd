package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGinkgoFd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinkgoFd Suite")
}
