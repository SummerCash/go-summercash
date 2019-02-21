# SummerCash

Go implementation of the SummerCash protocol.

[![Godoc Reference](https://img.shields.io/badge/godoc-reference-%23516aa0.svg)](https://godoc.org/github.com/SummerCash/go-summercash)
[![Go Report Card](https://goreportcard.com/badge/github.com/summercash/go-summercash)](https://goreportcard.com/report/github.com/summercash/go-summercash)
[![Build Status](https://travis-ci.com/SummerCash/go-summercash.svg?branch=master)](https://travis-ci.com/SummerCash/go-summercash)
[![Gluten Status](https://img.shields.io/badge/gluten-free-brightgreen.svg)](https://img.shields.io/badge/gluten-free-brightgreen.svg)

## Installation

### Getting the Source

```BASH
go get -u github.com/SummerCash/go-summercash
```

To install the go-summercash node for use in the Go bin, run:

```BASH
go get -u github.com/SummerCash/go-summercash && go install github.com/SummerCash/go-summercash
```

### Pre-compiled Binaries

For pre-compiled executable binaries, check the [releases](https://github.com/SummerCash/go-summercash/releases/latest) page.

## Usage

### Running the Node Daemon From Source

#### Running an Archival Node (Recommended)

```BASH
sudo go run main.go --archival
```

or, from the Go bin (same steps apply for all other bin commands):

```BASH
sudo go-summercash --archival
```

#### Running a Light Personal Node

```BASH
sudo go run main.go
```

#### Running the Node With Terminal Input

```BASH
sudo go run main.go --terminal
```

#### Connecting to a Running Node In Terminal Mode

```BASH
sudo go run main.go --terminal --rpc-address RPC-ADDRESS-HERE
```
