#!/usr/bin/env bash

git checkout origin/master -- c:/dev/OtusGolangPro/hw10_program_optimization/stats.go
go test -bench=BenchmarkGetDomainStat . -benchmem > old.txt

git checkout hw10_program_optimization -- c:/dev/OtusGolangPro/hw10_program_optimization/stats.go
go test -bench=BenchmarkGetDomainStat . -benchmem > new.txt

benchstat ./old.txt ./new.txt > ./diff.txt
