.PHONY: default help
define help

Makefile for DEWEY documentation
  setup   - install gitbook locally
  build   - build the documentation website
  clean   - remove the build
  watch   - run the webserver on port 4000 and watch for changes
  open    - Open the documentation website

endef
export help

help:
	@echo "$$help"

setup:
	npm install gitbook-cli

build:
	gitbook build

watch:
	gitbook serve

open:
	open http://localhost:4000

clean:


