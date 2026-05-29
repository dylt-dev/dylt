#! /usr/bin/env bash

main ()
{
	export ECO_DEPTH=2
	export ECO_GEN_TESTS=1
	export ECO_OUTPUT_FILE=./decode_stage1_test.go
	export ECO_TEST_COUNT=5

	rm eco/decode_genned_test.go
	rm eco/decode_stage1_test.go
	go test -test.fullpath=true -timeout 30s -run ^TestGenBootstrap$ github.com/dylt-dev/dylt/eco -v || return
	go test -test.fullpath=true -timeout 30s -run ^TestGenDecodeDeepTest$ github.com/dylt-dev/dylt/eco -v || return
	go test -test.fullpath=true -timeout 30s -run ^TestDecodeDeepType2 github.com/dylt-dev/dylt/eco -v || return
}

main
