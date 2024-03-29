# Solus makefile
run:
	@echo "Running..."
	@go run main.go
	@echo "Done."

generate_requirements:
	@echo "Generating requirements..."
	@make build
	@./solus.out requirements -f gen/messages.json -o gen/generated_requirements.yaml
	@make clean
	@echo "Done."

generate_code:
	@echo "Generating code..."
	@make build
	@./solus.out code -g $(shell pwd)/gen
	@make clean
	@echo "Done."

zip_result:
	@echo "Zipping result..."
	@zip -r gen.zip gen
	@echo "Done."

end_to_end:
	@echo "Running end to end..."
	@make run
	@make generate_requirements
	@make generate_code
	@make zip_result
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

lint:
	@echo "Linting..."
	@trunk fmt
	@echo "Done."

clean:
	@echo "Cleaning..."
	@rm -f solus.out
	@echo "Done."
