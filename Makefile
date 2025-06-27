SHELL := /bin/bash

default: install

build: *.go
	go build

completion: ./scripts/bash-completion/comp.sh 
	. ./scripts/bash-completion/comp.sh

install: build completion
	go install
	
run: build install
	dylt misc gen-etcd-run-script
		
