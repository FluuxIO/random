# Random

Fast Random Generator for use from single threaded data injection code. It is especially useful if you
need to load test a service while generating random data. It that case, the random generator locks can
become a bottleneck.

## How to use

You first need to create a `RandomUnsafe` generator:

```go
r := random.NewRandomUnsafe()
```

You will need to keep it as a reference. Please not that it is only safe to use from a single thread, but you can
generate several `RandomUnsafe` values at the beginning of different Goroutines.

You can then call one of the available method to generate specific values.

For example, to generate an Int:

```go
r.Intn(100)
```

To generate a variable length string between 15 and 20 characters:

```go
r.String(15, 25)
```

## Benchmark

Here are a few example showing the gain 

```bash
$ go test -bench .
goos: darwin
goarch: amd64
pkg: fluux.io/random
BenchmarkRandString-4            2000000               808 ns/op
BenchmarkRandomString-4         20000000                56.5 ns/op
BenchmarkRandId-4                1000000              1665 ns/op
BenchmarkRandomId-4              1000000              1053 ns/op
BenchmarkRandInt-4              50000000                32.0 ns/op
BenchmarkRandomInt-4            100000000               23.2 ns/op
PASS
ok      fluux.io/random 10.386s
```

