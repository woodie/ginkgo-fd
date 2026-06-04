# ginkgo-fd

The `ginko-fd` command uses [Ginkgo](https://github.com/onsi/ginkgo) under the hood to emulate the style of [RSpec](https://github.com/rspec/rspec) "format documentation" output.

## Installation

Install Ginkgo if you haven't already:

```
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

Then install ginkgo-fd:

```
go install github.com/woodie/ginkgo-fd@latest
```

Or build locally:

```
go mod tidy
go build -o ginkgo-fd
mv ginkgo-fd ~/go/bin/
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
    with a skipping report
      annotates the skipped spec
      prints the summary with skipped count
    when the report file is missing
      returns an error
  color output
    when not a TTY
      omits ANSI codes from passing leaf nodes
      omits ANSI codes from the summary
    when a TTY
      colors passing leaf nodes green
      colors the passing summary green
      colors failed leaf nodes red
      colors the failing summary red
      colors skipped leaf nodes cyan
  main routing
    when a .json argument is given
      formats the report file directly
    when runGinkgo writes a report
      uses a path outside the project directory
```
