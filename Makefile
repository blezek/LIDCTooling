export GOPATH:=$(shell pwd)
# Build version
VERSION:=$(shell git describe --always)-$(shell date +%F-%H-%M)
DIR:=$(shell pwd)

define help

Makefile for LIDCFetch
  build   - build the LIDCFetch app
  test    - run the tests
  vet     - look for trouble in the code...
  segment - Grab some files and segment

Example:
make build
bin/LIDCFetch fetch series 1.3.6.1.4.1.14519.5.2.1.6279.6001.303494235102183795724852353824
endef
export help

help:
	@echo "$$help"

deps:
	go get -d -v LIDCFetch/...

fmt:
	go fmt LIDCFetch/...

vet:
	go vet LIDCFetch/...

doc:
	godoc -http=:6060 -goroot=../go

build: bin/LIDCFetch

bin/LIDCFetch: deps
	go install -ldflags "-X main.Version=${VERSION}" LIDCFetch

dicom: build
	bin/LIDCFetch fetch image --extract dicom/ 1.3.6.1.4.1.14519.5.2.1.6279.6001.303494235102183795724852353824


segment: dicom
	mkdir -p dicom-segmented
	./gradlew jar
	(cd build/libs && java -jar LIDCTooling.jar segment ${DIR}/LIDC-XML-only/tcia-lidc-xml/157/158.xml ${DIR}/dicom ${DIR}/dicom-segmented)

test: build
	go test LIDCFetch/...
