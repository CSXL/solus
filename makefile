# Solus makefile
run:
	@echo "Running..."
	@go run main.go
	@echo "Done."

build:
	@echo "Building..."
	@go build -o solus.out
	@echo "Done."

install:
	@echo "Installing..."
	@go install
	@echo "Done."

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done."

clean:
	@echo "Cleaning..."
	@rm -f solus.out
	@echo "Done."
