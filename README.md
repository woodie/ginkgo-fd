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
ginkgo --json-report=out.json && ginkgo-fd out.json
```
