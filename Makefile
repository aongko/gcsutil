PACKAGE = gcsutil
SOURCES = $(shell find . -name '*.go' -type f -not -path './vendor/*')

vendor: glide.lock
	glide install --strip-vendor

.PHONY: install
install: vendor $(SOURCES)
	go install
