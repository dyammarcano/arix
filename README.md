# Nomad (dev purpose)

## Description

This tool is meant to be used for development purposes. It is a simple tool that allows you to run a variety of go tools in one place. It is meant to be used in a development environment and not in production.


## What is included

- [x] Go installer
- [x] Go tools
- [x] Go linters
- [x] Go formatters
- [x] Go vet
- [x] Go vet
- [x] Go test
- [x] Task
- [x] Goreleaser
- [x] Go vulncheck
- [x] GolangCI-Lint
- [x] Go coverage


## How to use

### Install

```bash
go install github.com/odacremolbap/arix@latest
```

### Run Check

This will run all the tools in the order they are listed above.

```bash
arix check
```

### Run Specific Tool

You can run a specific tool by using the following command:

```bash
arix run <tool>
```

For example, to run the linter you can use the following command:

```bash
arix run lint
```

### Run Specific Tool with Arguments

You can run a specific tool with arguments by using the following command:

```bash
arix run <tool> <args>
```

### Run Go install

You can run the go install command by using the following command:

```bash
arix run install
```

### Run Go test

You can run the go test command by using the following command:

```bash
arix run test
```

### Run Go install specific version

You can run the go install command with a specific version by using the following command:

```bash
arix run install <version>
```
