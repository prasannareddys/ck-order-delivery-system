help::
	@echo "======Available commands==========="
	@echo "make build"
	@echo "make test"
	@echo "make start"

test:
	@echo "====================="
	@echo "Running unit tests"
	@echo "====================="
	go test -race ./pkg/...

build:
	@echo "====================="
	@echo "Building Project"
	@echo "====================="
	go build -o ck-order-delivery-system
	@echo "Build complete"

start: build
	@echo "====================="
	@echo "running project"
	@echo "====================="
	./ck-order-delivery-system server --ops 2 --order-file-path ./data/orders.json
