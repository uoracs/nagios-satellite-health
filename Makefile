PHONY: build run

all: build

build:
	go build -o bin/nagios-satellite-health-check main.go

run:
	go run main.go
