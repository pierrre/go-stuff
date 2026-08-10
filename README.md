# Go stuff

Collection of various Go code: experiments, utilities, and demos.

## Packages

| Path | Description |
| --- | --- |
| [`func-mock`](func-mock) | Mock functions using reflection (`MockList`, `MockCount`, `MockSerial`). |
| [`func-mock/funcmocktest`](func-mock/funcmocktest) | Test helpers for `func-mock` with automatic failure reporting. |
| [`lowpriority`](lowpriority) | Run code on a low-priority OS thread, with a worker pool. |
| [`passwordentropy`](passwordentropy) | Estimate password entropy in bits from character classes. |
| [`fibonacci`](fibonacci) | Iterative Fibonacci string generator in 5 progressive variants (`v1`–`v5`). |
| [`bench-loop`](bench-loop) | Benchmarks comparing `b.Loop()` vs `b.N` iteration patterns. |
| [`reflect_pretty_assertauto`](reflect_pretty_assertauto) | Demos comparing `reflect`, `fmt`, [`pretty`](https://pkg.go.dev/github.com/pierrre/pretty), and [`assertauto`](https://pkg.go.dev/github.com/pierrre/assert/assertauto) output. |
| [`bug-proto-testify`](bug-proto-testify) | Reproduction of a `testify`/`proto.Equal` comparison bug. |
