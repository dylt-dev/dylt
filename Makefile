default: install

build: *.go
	go build

install: build
	go install
	
run: build install
	dylt misc gen-etcd-run-script
		
