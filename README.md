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

Then call `risky-func coverage.out`

## Example

```
+-------------------------+-----------+-----------------+
|          FILE           | FUNCTION  | UNCOVERED LINES |
+-------------------------+-----------+-----------------+
| ...y-func/risky_func.go | main      |              33 |
| ...y-func/risky_func.go | coverage  |              12 |
| ...y-func/risky_func.go | findFuncs |               7 |
| ...y-func/risky_func.go | Visit     |               6 |
| ...y-func/risky_func.go | findFile  |               5 |
| ...y-func/risky_func.go | Print     |               1 |
+-------------------------+-----------+-----------------+
```