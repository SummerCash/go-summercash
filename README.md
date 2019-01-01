# SummerCash

Go implementation of the SummerCash node

[![asciicast](https://asciinema.org/a/WyyI9GF7ycacyVm6x3G5IBtAk.svg)](https://asciinema.org/a/WyyI9GF7ycacyVm6x3G5IBtAk)

## Installation

To install the go-summercash node for direct source use, run:

```BASH
go get -u github.com/SummerCash/go-summercash
```

To install the go-summercash node for use in the Go bin, run:

```BASH
go get -u github.com/SummerCash/go-summercash && go install github.com/SummerCash/go-summercash
```

## Usage

### Running the Node Daemon From Source

#### Running an Archival Node (Recommended)

```BASH
sudo go run main.go --archival
```

#### Running a Light Personal Node

```BASH
sudo go run main.go
```

or, from the Go bin:

```BASH
sudo go-summercash
```

#### Running the Node With Terminal Input

```BASH
sudo go run main.go --terminal
```

#### Connecting to a Running Node In Terminal Mode

```BASH
sudo go run main.go --terminal --rpc-address RPC-ADDRESS-HERE
```