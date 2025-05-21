.PHONY: build
build:
	go build \
	-o application \
	cmd/*.go

all: build