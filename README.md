# ginkgo-fd

Formats Ginkgo JSON reports in the style of RSpec's `--format documentation` output.

## Installation

```bash
go install github.com/woodie/ginkgo-fd@latest
```

Or build locally:

```bash
go build -o ginkgo-fd
sudo mv ginkgo-fd /usr/local/bin/
```

## Usage

```bash
ginkgo --json-report=report.json && ginkgo-fd
```

Or with a custom path:

```bash
ginkgo --json-report=report.json && ginkgo-fd
```

Sample output:

```
Running Suite: GinkgoFd Suite - /Users/woodie/workspace/ginkgo-fd
=================================================================
Random Seed: 1780386104

Will run 10 of 10 specs
••••••••••

Ran 10 of 10 Specs in 0.007 seconds
SUCCESS! -- 10 Passed | 0 Failed | 0 Pending | 0 Skipped
PASS

Ginkgo ran 1 suite in 2.724026083s
Test Suite Passed

  GinkgoFd
    run
      with a passing report
        prints the suite name
        indents container hierarchy
        indents leaf nodes
        deduplicates shared hierarchy
        prints the summary
      with a failing report
        annotates the failed spec
        prints the failures section
        prints the failed examples list
        prints the summary with failure count
      when the report file is missing
        returns an error

Finished in 0.0067 seconds
10 examples, 0 failures
```
