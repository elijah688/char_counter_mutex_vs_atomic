.PHONY: test bench

test:
	go test -v -race -failfast

bench:
	go test -bench=. -benchtime=10s -timeout 20m -v ./...

