#!/usr/bin/env bash

OS=$(uname -s)

if [ "$OS" = "Linux" ]; then
    BASE_PATH="/home/void/dev"
else
    BASE_PATH="c:/dev"
fi

STATS_GO_PATH="$BASE_PATH/OtusGolangPro/hw10_program_optimization/stats.go"

git checkout origin/master -- "$STATS_GO_PATH"
go test -bench=BenchmarkGetDomainStat . -benchmem > old.txt

git checkout hw10_program_optimization -- "$STATS_GO_PATH"
go test -bench=BenchmarkGetDomainStat . -benchmem > new.txt

benchstat ./old.txt ./new.txt > ./diff.txt
