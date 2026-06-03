package main

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var passingReport = `[{
  "SuiteName": "Something",
  "RunTime": 15000000,
  "SpecReports": [
    {
      "ContainerHierarchyTexts": ["checkAttachmentDir", "when the path is missing"],
      "LeafNodeText": "creates the directory",
      "LeafNodeType": "It",
      "State": "passed"
    },
    {
      "ContainerHierarchyTexts": ["checkAttachmentDir", "when the path is missing"],
      "LeafNodeText": "does not error",
      "LeafNodeType": "It",
      "State": "passed"
    },
    {
      "ContainerHierarchyTexts": ["checkAttachmentDir", "when the path is a symlink"],
      "LeafNodeText": "does not error",
      "LeafNodeType": "It",
      "State": "passed"
    }
  ]
}]`

var failingReport = `[{
  "SuiteName": "Something",
  "RunTime": 15000000,
  "SpecReports": [
    {
      "ContainerHierarchyTexts": ["checkAttachmentDir", "when the path is missing"],
      "LeafNodeText": "creates the directory",
      "LeafNodeType": "It",
      "State": "failed",
      "Failure": {
        "Message": "Expected file to exist",
        "Location": {"FileName": "main_test.go", "LineNumber": 42}
      }
    }
  ]
}]`

var _ = Describe("GinkgoFd", func() {
	runReport := func(raw string) string {
		path := writeTempReport(raw)
		var buf strings.Builder
		Expect(run(path, &buf)).To(Succeed())
		return buf.String()
	}

	Describe("run", func() {
		var output string

		Context("with a passing report", func() {
			BeforeEach(func() { output = runReport(passingReport) })

			It("prints the suite name", func() {
				Expect(output).To(ContainSubstring("Something"))
			})

			It("indents container hierarchy", func() {
				Expect(output).To(ContainSubstring("checkAttachmentDir\n  when the path is missing"))
			})

			It("indents leaf nodes", func() {
				Expect(output).To(ContainSubstring("creates the directory"))
			})

			It("deduplicates shared hierarchy", func() {
				Expect(strings.Count(output, "when the path is missing")).To(Equal(1))
			})

			It("prints the summary", func() {
				Expect(output).To(ContainSubstring("3 examples, 0 failures"))
			})
		})

		Context("with a failing report", func() {
			BeforeEach(func() { output = runReport(failingReport) })

			It("annotates the failed spec", func() {
				Expect(output).To(ContainSubstring("creates the directory (FAILED - 1)"))
			})

			It("prints the failures section", func() {
				Expect(output).To(ContainSubstring("Failures:"))
				Expect(output).To(ContainSubstring("Expected file to exist"))
				Expect(output).To(ContainSubstring("main_test.go:42"))
			})

			It("prints the failed examples list", func() {
				Expect(output).To(ContainSubstring("Failed examples:"))
			})

			It("prints the summary with failure count", func() {
				Expect(output).To(ContainSubstring("1 examples, 1 failure"))
			})
		})

		Context("when the report file is missing", func() {
			It("returns an error", func() {
				var buf strings.Builder
				Expect(run("/nonexistent/report.json", &buf)).To(HaveOccurred())
			})
		})
	})

	Describe("color output", func() {
		Context("when not a TTY", func() {
			BeforeEach(func() { isTTY = false })
			AfterEach(func() { isTTY = false })

			It("omits ANSI codes from passing leaf nodes", func() {
				Expect(runReport(passingReport)).NotTo(ContainSubstring("\033["))
			})

			It("omits ANSI codes from the summary", func() {
				Expect(runReport(failingReport)).NotTo(ContainSubstring("\033["))
			})
		})

		Context("when a TTY", func() {
			BeforeEach(func() { isTTY = true })
			AfterEach(func() { isTTY = false })

			It("colors passing leaf nodes green", func() {
				Expect(runReport(passingReport)).To(ContainSubstring(green))
			})

			It("colors the passing summary green", func() {
				Expect(runReport(passingReport)).To(ContainSubstring(green))
			})

			It("colors failed leaf nodes red", func() {
				Expect(runReport(failingReport)).To(ContainSubstring(red))
			})

			It("colors the failing summary red", func() {
				Expect(runReport(failingReport)).To(ContainSubstring(red))
			})
		})
	})

	Describe("main routing", func() {
		Context("when a .json argument is given", func() {
			It("formats the report file directly", func() {
				path := writeTempReport(passingReport)
				var buf strings.Builder
				Expect(run(path, &buf)).To(Succeed())
				Expect(buf.String()).To(ContainSubstring("Something"))
			})
		})

		Context("when runGinkgo writes a report", func() {
			It("uses a path outside the project directory", func() {
				Expect(ginkgoReportPath()).To(HavePrefix(os.TempDir()))
			})
		})
	})
})
