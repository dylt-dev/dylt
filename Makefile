SHELL := /bin/bash

default: install

build: *.go
	go build

completion: ./scripts/bash-completion/_daylight.sh 
	. ./scripts/bash-completion/_daylight.sh

cross: darwin linux windows

darwin: *.go
	GOOS=darwin go build

install: build completion
	go install

linux: *.go
	GOOS=linux go build
	
run: build install
	dylt misc gen-etcd-run-script

windows: *.go
	GOOS=windows go build
		
