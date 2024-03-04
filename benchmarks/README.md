# Benchmarking

[`benchmark.go`](./benchmark.go) offers a way to time the execution phase of a program.
All timings provided are in **microseconds** (µs).

## How to evaluate (runtime) performance

Use the `--benchmark` flag to evaluate performance.
Optional flags include `--maxcores <number of cores>` (sets GOMAXPROCS) and `--repeat <number of times>` (to repeat a run multiple times) for fine-tuning tests.

```bash
./grits --benchmark examples/nat_double.grits
./grits --benchmark --maxcores 1 --repeat 5 examples/nat_double.grits
```

To run all pre-configured benchmarks: `./grits --sample-benchmarks`.

## How to interpret the output

When `--benchmark` (or `--sample-benchmarks`) is used, the program is timed using both the non-polarized (synchronous) transition semantics (v1) and the polarized version (v2).
The number of processes spawned is measured as well.

Results are collected in the `benchmark-results` directory.
Two files are outputted for each result: `file-benchmark-X.csv` and `file-benchmark-detailed-X.csv` (where `X` is the number of cores used).
The latter provides each result for each run, whilst the form only shows the averaged result.
These csv files contain the following columns:

- *name*: name of file being checked
- *timeSyncV1NP*: time taken to evaluate file (using v1)
- *processCountSyncV1NP*: number of processes spawn (when using v1)
- *timeAsyncV2*: time taken to evaluate file (using v2-async)
- *processCountAsyncV2*: number of processes spawn (when using v2-async)
<!-- - *timeSyncV2*: time taken to evaluate file (using v2-sync)
- *processCountSyncV2*: number of processes spawn (when using v2-sync) -->

## Sample results

Sample benchmarks for [`nat-double-1.go`](./compare/nat-double/nat-double-1.grits) can be found in [`nat-benchmarks-10.csv`](./compare/nat-double/sample-results/nat-benchmarks-10.csv) and [`nat-detailed-benchmarks-10.csv`](./compare/nat-double/sample-results/nat-detailed-benchmarks-10.csv).
