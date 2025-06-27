SHELL := /bin/bash

default: install

build: *.go
	go build

completion: ./scripts/bash-completion/_daylight.sh 
	. ./scripts/bash-completion/_daylight.sh

install: build completion
	go install
	
run: build install
	dylt misc gen-etcd-run-script
		
