# Risky Func

Drive higher confidence in making changes by detecting large blocks of untested functionality.

## Background

There are large code bases with low test coverages and there is low confidence to make changes in them. This low
confidence stems from the fact that when a change is made there is no guarantee that it will be correct from
a business perspective.

## Building

```
go build
```

## Usage

Go to the folder where `coverage.out` is generated via `go test ./... -coverprofile=coverage.out`

Run
```shell
./risky-func -f coverage.out
```
## Example

```
+-------------------------+------------------+-----------------+--------+
|          FILE           |     FUNCTION     | UNCOVERED LINES | IMPACT |
+-------------------------+------------------+-----------------+--------+
| ...y-func/risky_func.go | main             |              28 |   42.4 |
| ...y-func/risky_func.go | getFunctionInfos |               8 |   12.1 |
| ...y-func/risky_func.go | findFuncs        |               7 |   10.6 |
| ...y-func/risky_func.go | Visit            |               7 |   10.6 |
| ...y-func/risky_func.go | printTable       |               5 |    7.6 |
| ...y-func/risky_func.go | Print            |               5 |    7.6 |
| ...y-func/risky_func.go | trimString       |               1 |    1.5 |
| ...y-func/risky_func.go | coverage         |               1 |    1.5 |
| ...y-func/risky_func.go | findFile         |               1 |    1.5 |
+-------------------------+------------------+-----------------+--------+
```