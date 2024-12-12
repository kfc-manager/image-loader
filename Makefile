all: build run clean

build:
	@go build -o ./bin

run:
	@./bin

clean:
	@rm -rf ./bin
