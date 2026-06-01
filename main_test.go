package main

import (
    "strings"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var passingReport = `[{
    "SuiteName": "Lambada",
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
    "SuiteName": "Lambada",
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
        Context("with a passing report", func() {
            var output string
            BeforeEach(func() { output = runReport(passingReport) })

            It("prints the suite name", func() {
                Expect(output).To(ContainSubstring("Lambada"))
            })
            It("indents container hierarchy", func() {
                Expect(output).To(ContainSubstring("  checkAttachmentDir"))
                Expect(output).To(ContainSubstring("    when the path is missing"))
            })
            It("indents leaf nodes", func() {
                Expect(output).To(ContainSubstring("      creates the directory"))
            })
            It("deduplicates shared hierarchy", func() {
                Expect(strings.Count(output, "when the path is missing")).To(Equal(1))
            })
            It("prints the summary", func() {
                Expect(output).To(ContainSubstring("3 examples, 0 failures"))
            })
        })

        Context("with a failing report", func() {
            var output string
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
})
