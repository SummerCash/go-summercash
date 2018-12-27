# SummerCash

Go implementation of the SummerTech digital currency

## Specification

### Data Structure

Unlike most other cryptographically secured digital currencies, the SummerCash blockchain does not make use of a single central blockchain, but a collective of multiple interlacing directed acyclic graphs belonging to the users of the aforementioned currency. Futhermore, the SummerCash blockchain makes use of an overarching distributed "coordination chain" containing lookup metadata regarding each of the interlacing sub-chains.

## Usage

### Running the Node Daemon

```BASH
sudo go run main.go
```

### Running the Node With Terminal Input

```BASH
sudo go run main.go --terminal
```