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
Ginkgo ran 1 suite in 2.710928042s
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

Finished in 0.0060 seconds
10 examples, 0 failures
```
