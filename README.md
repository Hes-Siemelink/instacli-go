# Instacli

Instacli is a command-line tool for creating and running CLI applications with light scripting.

## Installation

```bash
go install instacli/cmd/cli@latest
```

## Usage

```bash
cli [global options] file | directory [command options]
```

### Global Options

- `--help, -h`: Print help on a script or directory and does not run anything
- `--output, -o`: Print the output at the end of the script in Yaml format
- `--output-json, -j`: Print the output at the end of the script in Json format
- `--non-interactive, -q`: Indicate that Instacli should not prompt for user input
- `--debug, -d`: Run in debug mode. Prints stacktraces when an error occurs.

## Development

### Building

```bash
go build -o cli cmd/cli/main.go
```

### Running Tests

```bash
go test ./...
```

## Project Structure

- `cmd/cli`: Main command-line interface
- `pkg/cli`: Core functionality for script execution

## License

MIT

# Instacli Go Implementation

## Build and Test Instructions

### Prerequisites
- Go 1.16 or later
- The `instacli/` directory is **read-only** and contains the Instacli specification (`instacli/instacli-spec`) and the Kotlin reference implementation (`instacli/src`).
+ Do not modify anything inside `instacli/`.

### Setting up the Spec Directory

Set the `INSTACLI_SPEC` environment variable to the path of your local Instacli spec directory. For example:

```sh
export INSTACLI_SPEC=instacli/instacli-spec
```

This allows the Go tests to locate and use the Instacli spec and test files.

### Building

```sh
go build ./...
```

### Running Tests

```sh
export INSTACLI_SPEC=instacli/instacli-spec
# or the absolute path to your instacli/instacli-spec directory

go test ./...
```

This will run all unit and integration tests, including those that use the Instacli spec and test files.

### Notes
- The test driver expects the spec to be at the path specified by `INSTACLI_SPEC` (e.g., `instacli/instacli-spec`).
- If you update the spec, you do not need to rebuild the Go binary, but you do need to re-run the tests.
- For CI, ensure the environment variable is set before running `go build` or `go test`. 