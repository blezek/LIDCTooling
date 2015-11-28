.PHONY: default help
define help

Makefile for DEWEY documentation
  setup   - install hugo, styles, gulp locally
  build   - build the documentation website
  clean   - remove the build
  watch   - run the webserver on port 1313 and watch for changes
  stop    - stop the Hugo server
  restart - restart the Hugo server
  open    - Open the documentation website

endef
export help

help:
	@echo "$$help"

setup:
	go get -u -v github.com/spf13/hugo
	npm install
	bower install

build:
	gulp
	hugo

watch: hugo.pid

restart:
	${MAKE} stop
	${MAKE} watch

hugo.pid:
	{ hugo serve & echo $$! > $@; }

stop: hugo.pid
	kill `cat $<` && rm $<

open:
	open http://localhost:1313

clean:
	rm -rf public

