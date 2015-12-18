.PHONY: help setup build npm default

default: help

define help

Makefile for MDS POC documentat
  setup   - install gitbook, assumes Node / npm are installed
  build   - build the documentation, in documentation/static
  watch   - watch for documentation changes, serves on http://localhost:4000
  open    - Open the documentation website (http://localhost:4000)

endef
export help

build: setup
	node_modules/.bin/gitbook build

help:
	@echo "$$help"

open:
	open http://localhost:4000

watch: gitbook.pid setup

restart:
	${MAKE} stop
	${MAKE} watch

gitbook.pid:
	{ node_modules/.bin/gitbook serve & echo $$! > $@; }

stop: 
	test -e gitbook.pid && kill `cat gitbook.pid` && rm gitbook.pid

publish: gh-pages build
	rsync -ar _book/ gh-pages
	(cd gh-pages && git add . && git commit -m "publish" && git push)

gh-pages:
	git clone `git remote -v | grep origin | grep fetch | awk '{ print $$2 }'` gh-pages
	(cd gh-pages && git checkout gh-pages)

clean:
	rm -rf _book


## Setup to build documentation
setup: node_modules/gitbook-cli

# installs gitbook
node_modules/gitbook-cli:
	@[ -x "`which npm 2>/dev/null`" ] || (printf "\n=====\nCould not find npm in your PATH, please install from http://nodejs.org\n=====\n\n"; exit 1;)
	npm install gitbook-cli

