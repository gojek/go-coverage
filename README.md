<br />
<p align="center"><h3 align="center">go-coverage</h3>

  <p align="center">
    Increase code coverage of Go projects
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary><h2 style="display: inline-block">Table of Contents</h2></summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

The key challenge with large code bases with low test coverage is to prioritize which sections of code to test first.

The standard coverage tools tell about the code coverage percentage and what is covered and uncovered however it doesn't give an input on which functions to cover first and what will be the impact of covering them.

This tool addresses the challenge by providing the sorted list of functions to cover and the impact associated with covering it.

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

You'll need Go installed to use this tool. [Here](https://golang.org/doc/install) is the installation instructions for Go.

### Installation

Via go get
```shell
go get -u github.com/gojekfarm/go-coverage
```

## Usage

### Prerequisites

Generate the coverage profile for your Go codebase, usually done via
```shell
go test ./... -coverprofile=coverage.out
```

### Get lines uncovered greater than 10

```shell
go-coverage -f coverage.out --line-filter 10
```

### Get trimmed file names

```shell
go-coverage -f coverage.out --line-filter 10 --trim
```

```shell
+-------------------------+-------------------------------------+-----------------+--------+
|          FILE           |              FUNCTION               | UNCOVERED LINES | IMPACT |
+-------------------------+-------------------------------------+-----------------+--------+
| ...ice/config/config.go | RadiusForClosestDriverByServicetype |              26 |    1.9 |
| ...ice/config/config.go | RadiusForServicetype                |              26 |    1.9 |
| ...ice/config/config.go | AliceDriverLimit                    |              26 |    1.9 |
| ...ice/config/config.go | ConsumerDriverLimitByServicetype    |              26 |    1.9 |
| .../service/handlers.go | findDriver                          |              19 |    1.4 |
| ...ice/extern/driver.go | driverAllocationStatusFromAPI       |              19 |    1.4 |
| .../service/handlers.go | updateDriverVehicleTags             |              18 |    1.3 |
| ...ice/config/config.go | ConsumerDriverLimit                 |              14 |    1.0 |
| ...vice/service/cron.go | startCrons                          |              14 |    1.0 |
| ...ice/config/config.go | RadiusForVehicleType                |              13 |    0.9 |
| ...ice/config/config.go | matchVehicleType                    |              12 |    0.9 |
| ...rvice/service/api.go | startServer                         |              11 |    0.8 |
+-------------------------+-------------------------------------+-----------------+--------+
```

### Exclude file name pattern

```shell
go-coverage -f coverage.out --exclude ".*config.*" --line-filter 10 --trim
```

```shell
+-------------------------+-------------------------------+-----------------+--------+
|          FILE           |           FUNCTION            | UNCOVERED LINES | IMPACT |
+-------------------------+-------------------------------+-----------------+--------+
| .../service/handlers.go | findDriver                    |              19 |    1.4 |
| ...ice/extern/driver.go | driverAllocationStatusFromAPI |              19 |    1.4 |
| .../service/handlers.go | updateDriverVehicleTags       |              18 |    1.3 |
| ...vice/service/cron.go | startCrons                    |              14 |    1.0 |
| ...rvice/service/api.go | startServer                   |              11 |    0.8 |
+-------------------------+-------------------------------+-----------------+--------+
```

## Roadmap

- [ ] Support generation of HTML
- [ ] Integrate with gitlab
