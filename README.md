# ginkgo-fd

Formats [Ginkgo](https://github.com/onsi/ginkgo) JSON reports in the style of RSpec's `--format documentation` output, with color support.

## Installation

```
go install github.com/woodie/ginkgo-fd@latest
```

Or build locally:

```
go build -o ginkgo-fd
sudo mv ginkgo-fd /usr/local/bin/
```

## Usage

Run as a wrapper around `ginkgo`, passing any arguments through:

```
ginkgo-fd
ginkgo-fd ./...
ginkgo-fd -v ./mypackage
```

Or format an existing report file directly:

```
ginkgo-fd report.json
```

Sample output:

```
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

Finished in 0.0043 seconds
10 examples, 0 failures
```
