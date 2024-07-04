PHONY: benchmark
benchmark:
	@echo "Benchmarking..."
	@redis-benchmark -p 6379 -t set,get -n 100000 -c 100 -q
	@echo "Benchmarking done."

PHONY: build
build:
	@echo "Building..."
	@go build -o bin/redigo ./main.go
	@echo "Building done."

PHONY: run
run: build
	@echo "Running..."
	@./bin/redigo
	@echo "Running done."