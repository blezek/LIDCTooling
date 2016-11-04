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
  io              - 1 r3.large for install / data movement
  memory          - 20x r3.large (large memory) spot instances

Example:
make build
bin/LIDCFetch fetch series 1.3.6.1.4.1.14519.5.2.1.6279.6001.303494235102183795724852353824
endef
export help

help:
	@echo "$$help"

build:
	./gradlew installDist

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
