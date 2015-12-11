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
	go get -u github.com/spf13/hugo
	npm install
	bower install

build:
	gulp
	hugo

watch: hugo.pid gulp.pid

restart:
	${MAKE} stop
	${MAKE} watch

hugo.pid:
	{ hugo serve & echo $$! > $@; }

gulp.pid:
	{ gulp watch & echo $$! > $@; }

stop: 
	test -e hugo.pid && kill `cat hugo.pid` && rm hugo.pid
	test -e gulp.pid && kill `cat gulp.pid` && rm gulp.pid

open:
	open http://localhost:1313

publish: gh-pages build
	rsync -ar public/ gh-pages
	(cd gh-pages && git add . && git commit -m "publish" && git push)

gh-pages:
	git clone `git remote -v | grep origin | grep fetch | awk '{ print $$2 }'` gh-pages
	(cd gh-pages && git checkout gh-pages)


clean:
	rm -rf public

