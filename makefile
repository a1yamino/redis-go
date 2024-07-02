PHONY: benchmark
benchmark:
	@echo "Benchmarking..."
	@redis-benchmark -p 6379 -t set,get -n 100000 -c 100 -q
	@echo "Benchmarking done."