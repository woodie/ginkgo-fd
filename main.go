package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Minimal types matching Ginkgo's JSON report schema.

type Report struct {
	SuiteName   string       `json:"SuiteName"`
	RunTime     float64      `json:"RunTime"`
	SpecReports []SpecReport `json:"SpecReports"`
}

type SpecReport struct {
	ContainerHierarchyTexts []string  `json:"ContainerHierarchyTexts"`
	LeafNodeText            string    `json:"LeafNodeText"`
	LeafNodeType            string    `json:"LeafNodeType"`
	State                   string    `json:"State"`
	RunTime                 float64   `json:"RunTime"`
	Failure                 *Failure  `json:"Failure,omitempty"`
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

func run(reportPath string) error {
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %w", reportPath, err)
	}

	// Ginkgo writes an array of Report (one per suite).
	var reports []Report
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
		fmt.Println(report.SuiteName)
		totalRunTime += report.RunTime

		// Track the last printed hierarchy to avoid repeating shared prefixes.
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

			// Print any new levels of container hierarchy.
			hierarchy := spec.ContainerHierarchyTexts
			divergeAt := 0
			for divergeAt < len(prevHierarchy) && divergeAt < len(hierarchy) &&
				prevHierarchy[divergeAt] == hierarchy[divergeAt] {
				divergeAt++
			}
			for i := divergeAt; i < len(hierarchy); i++ {
				fmt.Printf("%s%s\n", strings.Repeat("  ", i+1), hierarchy[i])
			}

			// Print the It text with state annotation.
			depth := len(hierarchy) + 1
			indent := strings.Repeat("  ", depth)
			label := spec.LeafNodeText
			switch spec.State {
			case "failed", "panicked":
				n := len(failures) + 1
				label = fmt.Sprintf("%s (FAILED - %d)", label, n)
				failures = append(failures, failureEntry{
					n:        n,
					full:     append(append([]string{report.SuiteName}, hierarchy...), spec.LeafNodeText),
					message:  spec.Failure.Message,
					location: fmt.Sprintf("%s:%d", spec.Failure.Location.FileName, spec.Failure.Location.LineNumber),
				})
			case "pending":
				label = fmt.Sprintf("%s (PENDING)", label)
			case "skipped":
				label = fmt.Sprintf("%s (SKIPPED)", label)
			}
			fmt.Printf("%s%s\n", indent, label)

			prevHierarchy = hierarchy
		}
		fmt.Println()
	}

	// Failures section.
	if len(failures) > 0 {
		fmt.Println("Failures:")
		for _, f := range failures {
			fmt.Printf("\n  %d) %s\n", f.n, strings.Join(f.full, " "))
			for _, line := range strings.Split(strings.TrimSpace(f.message), "\n") {
				fmt.Printf("     %s\n", line)
			}
			fmt.Printf("     # %s\n", f.location)
		}
		fmt.Println()
	}

	// Summary line.
	fmt.Printf("Finished in %s\n", formatDuration(totalRunTime))
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
	fmt.Println(strings.Join(parts, ", "))

	// Failed examples list.
	if len(failures) > 0 {
		fmt.Println("\nFailed examples:")
		for _, f := range failures {
			fmt.Printf("\n  # %s\n", strings.Join(f.full, " "))
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

func main() {
	path := "report.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	if err := run(path); err != nil {
		fmt.Fprintf(os.Stderr, "ginkgo-fd: %v\n", err)
		os.Exit(1)
	}
}
