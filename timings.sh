#!/bin/bash
# A sample Bash script, by Ryan
echo Hello World!


go run . --benchmark --maxcores 1 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 2 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 3 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 4 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 5 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 6 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 7 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 8 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 9 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi
go run . --benchmark --maxcores 10 --repeat 30 ./benchmarks/compare/nat-double/nat-double-8.phi