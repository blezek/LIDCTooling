export GOPATH:=$(shell pwd)
# Build version
VERSION:=$(shell git describe --always)-$(shell date +%F-%H-%M)
DIR:=$(shell pwd)
CLUSTER ?= lidc
TEMPLATE ?= smallcluster

define help

Makefile for LIDCFetch
  build   - build the LIDCFetch app
  test    - run the tests
  vet     - look for trouble in the code...
  segment - Grab some files and segment

Cluster (make CLUSTER=test TEMPLATE=smallcluster)
  cluster-start      - start up starcluster
  cluster-install    - copy ClusterSoftware to /software
  cluster-ssh        - ssh (as sgemaster) to cluster
  cluster-terminate  - terminate (destroy) starcluster

Cluster Templates
  smallcluster    - testing, 2 micro nodes
  lidc            - production, 10 c3.xlarge spot instances
  io              - 1 c3.large for install / data movement

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

evaluate: bin/LIDCFetch
	rm run_test.db
	bin/LIDCFetch --verbose evaluate --db run_test.db segmented/*

run: build bin/LIDCFetch
	rm -rf segmented/1.3.6.1.4.1.14519.5.2.1.6279.6001.303494235102183795724852353824/
	env PATH=../ChestImagingPlatform/build/CIP-build/bin/:python/:${PATH} \
	bin/LIDCFetch \
              process \
              --extract build/install/LIDCTooling/bin/Extract \
              --fetch bin/LIDCFetch \
              --evaluate python/evaluateSegmentation.py \
              --dicom dicom/ \
              --segmented segmented/ \
              ClusterSoftware/tcia-lidc-xml/157/158.xml
	bin/LIDCFetch \
	      evaluate \
	      --db test_run.db \
              segmented/*

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
