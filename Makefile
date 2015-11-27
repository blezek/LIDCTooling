export GOPATH:=$(shell pwd)
# Build version
VERSION:=$(shell git describe --always)-$(shell date +%F-%H-%M)
DIR:=$(shell pwd)
CLUSTER ?= test
TEMPLATE ?= smallcluster

define help

Makefile for LIDCFetch
  build   - build the LIDCFetch app
  test    - run the tests
  vet     - look for trouble in the code...
  segment - Grab some files and segment

Cluster (make CLUSTER=test)
  cluster-start      - start up starcluster
  cluster-terminate  - terminate (destroy) starcluster
  cluster-install    - copy ClusterSoftware to /software

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
	./gradlew installDist

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

cluster-ssh:
	(source venv/bin/activate && starcluster sshmaster --user sgeadmin ${CLUSTER})

cluster-start:
	(source venv/bin/activate && starcluster start -c ${TEMPLATE} ${CLUSTER})

cluster-terminate:
	(source venv/bin/activate && starcluster terminate --confirm ${CLUSTER})

cluster-loadbalance:
	(source venv/bin/activate && starcluster loadbalance -m 10 -w 10 -a 4 ${CLUSTER})

cluster-install:
	devops/cluster-install.sh ${CLUSTER}
