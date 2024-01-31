# Benchmarking

[`benchmark.go`](./benchmark.go) offer a way to time the execution phase of a program.
All timings provided are in **microseconds** (µs).

## How to evaluate the performance

Use the `--benchmark` flag to evaluate performance.
Optional flags include `--maxcores <number of cores>` (sets GOMAXPROCS) and `--repeat <number of times>` (to repeat a run multiple times) for fine-tuning tests.

```bash
./phi --benchmark examples/nat_double.phi
./phi --benchmark --maxcores 1 --repeat 5 examples/nat_double.phi
```

To run all pre-configured benchmarks: `./phi --benchmarks`.

## How to interpret the output

When `--benchmark` or `--benchmarks` is used, the program is times using the non-polarized (synchronous) transition semantics (v1) and the polarized version (v2).
The number of processes spawned is measured as well.

Two files are outputted for each result: `file-benchmark-X.csv` and `file-benchmark-detailed-X.csv` (where `X` is the number of cores used).
The latter provides each result for each run, whilst the form only shows the averaged result.
These csv files contain the following columns:

- *name*: name of file being checked
- *timeNonPolarizedSync*: time taken to evaluate file (using v1)
- *processCountNonPolarizedSync: number of processes spawn (when using v1)
- *timeNormalAsync*: time taken to evaluate file (using v2-async)
- *processCountNormalAsync*: number of processes spawn (when using v2-async)
- *timeNormalSync*: time taken to evaluate file (using v2-sync)
- *processCountNormalSync*: number of processes spawn (when using v2-sync)
