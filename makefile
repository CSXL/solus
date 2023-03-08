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

clean:
	@echo "Cleaning..."
	@rm -f solus
	@echo "Done."
