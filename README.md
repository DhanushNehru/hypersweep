# HyperSweep ⚡

A blazing-fast, highly resilient, concurrent link checker written natively in Go. 
HyperSweep rips through your local files, extracts URLs, and validates them across the inter-webs concurrently using Go's lightweight worker-pool architecture. 

It was built as a zero-dependency, ultra-lightweight alternative to link checkers like Lychee.

## Features
- **Concurrent Engine:** Execute 50+ network requests simultaneously without overloading file descriptors or crashing your OS.
- **Smart Retries:** Capable of recovering gracefully from aggressive firewalls (retrying blocked `HEAD` requests as standard `GET` traffic).
- **Static Binary:** Single compiled binary. Extremely fast execution with minimal memory footprint.
- **Made for CI/CD:** Native Github Actions Docker integration. Fast. Safe. Unbroken chains.

## Usage as a GitHub Action

The absolute best way to use HyperSweep is to have it automatically check the links across your GitHub repository every time someone pushes new code.

Simply add this to a workflow file in your project (e.g. `.github/workflows/link-checker.yml`):

```yaml
name: Check Links
on: [push, pull_request]

jobs:
  hypersweep:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4
      
      - name: Validate Links using HyperSweep
        uses: DhanushNehru/hypersweep@main
        with:
          path: '.'
          workers: '100'
          timeout: '15'
```

### Action Configuration (Inputs)
| Input | Description | Default |
|-------|-------------|---------|
| `path` | The path or directory where standard Markdown, Text, or HTML files exist. | `.` (Repo Root) |
| `workers` | The max number of simultaneous HTTP requests to execute. | `50` |
| `timeout` | Connection & read timeout (seconds) for each web request. | `10` |

### Installing the CLI Locally

If you have Go installed on your machine, you can run this purely as a CLI:

```bash
git clone https://github.com/DhanushNehru/hypersweep.git
cd hypersweep

go run cmd/hypersweep/main.go -path ./docs -workers 100
```
