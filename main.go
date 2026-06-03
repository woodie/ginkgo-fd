package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/term"
)

var isTTY = term.IsTerminal(int(os.Stdout.Fd()))

func colorize(code, s string) string {
	if !isTTY {
		return s
	}
	return code + s + "\033[0m"
}

const (
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

type SuiteReport struct {
	SuiteName   string       `json:"SuiteName"`
	RunTime     float64      `json:"RunTime"`
	SpecReports []SpecResult `json:"SpecReports"`
}

type SpecResult struct {
	ContainerHierarchyTexts []string `json:"ContainerHierarchyTexts"`
	LeafNodeText            string   `json:"LeafNodeText"`
	LeafNodeType            string   `json:"LeafNodeType"`
	State                   string   `json:"State"`
	RunTime                 float64  `json:"RunTime"`
	Failure                 *Failure `json:"Failure,omitempty"`
}

type Failure struct {
	Message  string   `json:"Message"`
	Location Location `json:"Location"`
}

type Location struct {
	FileName   string `json:"FileName"`
	LineNumber int    `json:"LineNumber"`
}

func formatDuration(ns float64) string {
	d := time.Duration(ns)
	if d < time.Second {
		return fmt.Sprintf("%.4f seconds", d.Seconds())
	}
	return fmt.Sprintf("%.2f seconds", d.Seconds())
}

func run(reportPath string, out io.Writer) error {
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", reportPath, err)
	}

	var reports []SuiteReport
	if err := json.Unmarshal(data, &reports); err != nil {
		return fmt.Errorf("cannot parse JSON: %w", err)
	}

	totalSpecs := 0
	totalFailed := 0
	totalPending := 0
	totalSkipped := 0
	var totalRunTime float64
	var failures []failureEntry

	for _, report := range reports {
		fmt.Fprintln(out, report.SuiteName)
		totalRunTime += report.RunTime

		var prevHierarchy []string

		for _, spec := range report.SpecReports {
			if spec.LeafNodeType != "It" {
				continue
			}

			totalSpecs++
			switch spec.State {
			case "failed", "panicked", "interrupted":
				totalFailed++
			case "pending":
				totalPending++
			case "skipped":
				totalSkipped++
			}

			hierarchy := spec.ContainerHierarchyTexts
			divergeAt := 0
			for divergeAt < len(prevHierarchy) && divergeAt < len(hierarchy) &&
				prevHierarchy[divergeAt] == hierarchy[divergeAt] {
				divergeAt++
			}

			for i := divergeAt; i < len(hierarchy); i++ {
				fmt.Fprintf(out, "%s%s\n", strings.Repeat("  ", i), hierarchy[i])
			}

			depth := len(hierarchy)
			indent := strings.Repeat("  ", depth)
			label := spec.LeafNodeText

			switch spec.State {
			case "failed", "panicked":
				n := len(failures) + 1
				label = fmt.Sprintf("%s (FAILED - %d)", label, n)
				label = colorize(red, label)
				failures = append(failures, failureEntry{
					n:        n,
					full:     append(append([]string{report.SuiteName}, hierarchy...), spec.LeafNodeText),
					message:  spec.Failure.Message,
					location: fmt.Sprintf("%s:%d", spec.Failure.Location.FileName, spec.Failure.Location.LineNumber),
				})
			case "pending":
				label = colorize(yellow, fmt.Sprintf("%s (PENDING)", label))
			case "skipped":
				label = colorize(yellow, fmt.Sprintf("%s (SKIPPED)", label))
			default:
				label = colorize(green, label)
			}

			fmt.Fprintf(out, "%s%s\n", indent, label)
			prevHierarchy = hierarchy
		}

		fmt.Fprintln(out)
	}

	if len(failures) > 0 {
		fmt.Fprintln(out, "Failures:")
		for _, f := range failures {
			fmt.Fprintf(out, "\n  %d) %s\n", f.n, strings.Join(f.full, " "))
			for _, line := range strings.Split(strings.TrimSpace(f.message), "\n") {
				fmt.Fprintf(out, "     %s\n", line)
			}
			fmt.Fprintf(out, "     # %s\n", f.location)
		}
		fmt.Fprintln(out)
	}

	fmt.Fprintf(out, "Finished in %s\n", formatDuration(totalRunTime))

	parts := []string{fmt.Sprintf("%d examples", totalSpecs)}
	if totalFailed > 0 {
		parts = append(parts, fmt.Sprintf("%d failure", totalFailed))
	} else {
		parts = append(parts, "0 failures")
	}
	if totalPending > 0 {
		parts = append(parts, fmt.Sprintf("%d pending", totalPending))
	}
	if totalSkipped > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", totalSkipped))
	}
	summary := strings.Join(parts, ", ")
	if totalFailed > 0 {
		summary = colorize(red, summary)
	} else {
		summary = colorize(green, summary)
	}
	fmt.Fprintln(out, summary)

	if len(failures) > 0 {
		fmt.Fprintln(out, "\nFailed examples:")
		for _, f := range failures {
			fmt.Fprintf(out, "\n  # %s\n", strings.Join(f.full, " "))
		}
	}

	return nil
}

type failureEntry struct {
	n        int
	full     []string
	message  string
	location string
}

func runGinkgo(args []string) int {
	reportPath := filepath.Join(os.TempDir(), "ginkgo-fd-report.json")
	defer os.Remove(reportPath)

	ginkgoArgs := append([]string{"--json-report=" + reportPath}, args...)
	cmd := exec.Command("ginkgo", ginkgoArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// ginkgo failed (e.g. test failures) — still format the report
			if _, statErr := os.Stat(reportPath); statErr == nil {
				fmt.Fprintln(os.Stdout)
				run(reportPath, os.Stdout)
			}
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "ginkgo-fd: %v\n", err)
		return 1
	}

	fmt.Fprintln(os.Stdout)
	if err := run(reportPath, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "ginkgo-fd: %v\n", err)
		return 1
	}
	return 0
}

func main() {
	args := os.Args[1:]

	// A single .json argument — format it directly.
	if len(args) == 1 && strings.HasSuffix(args[0], ".json") {
		if err := run(args[0], os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "ginkgo-fd: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Everything else (including no args) runs ginkgo as a wrapper.
	os.Exit(runGinkgo(args))
}
