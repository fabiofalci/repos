VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_DATE := $(shell date +"%s")

default: build

test:
	go test -v ./...

build:
	glide -q install
	go build -ldflags "\
		 -X main.version=$(VERSION) \
		 -X main.commit=$(COMMIT) \
		 -X main.buildDate=$(BUILD_DATE)"\
		 -o bundles/repos

clean:
	rm -rf bundles/

bundles:
	mkdir -p bundles/container

