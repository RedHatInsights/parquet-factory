---
layout: page
title: CI
nav_order: 3
---

# Continuous Integration

The Parquet Factory repository is configured to use the GitLab CI pipelines in
order to perform several checks and validations for every pushed commit:

* Unit tests using the standard tool `go test`.
* Several formatting and lintian checks, using `go fmt`, `go vet` and `golint`.
* `gocyclo` to report functions with too high cyclomatic complexity.
* `goconst` to search repeated strings in the source code.
* Security inspection of the code using `gosec` by scanning the Go AST.
* `ineffassign` to detect and print all ineffectual assignments in Go code.
* Detection of unchecked errors in the program with `errcheck`.
* `shellcheck` that analyze all the Shell scripts used in this repository.
* `abcgo` to measure ABC metrics and verify that the code is under an specified
  threshold.

All the checks above are executed for every commit and merge request. The
history of the pipelines execution can be found
[here](https://github.com/RedHatInsights/parquet-factory/-/pipelines).


