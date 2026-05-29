#! /usr/bin/env bash

main ()
{
	export ECO_GEN_TESTS=1
	export ECO_DEPTH=1000
	export ECO_OUTPUT_FILE=./decode_genned_stage1_test.go
	export ECO_OUTPUT_FILE2=./decode_genned_biguns_test.go
	export ECO_TEST_COUNT=1

	if [[ -f "eco/$ECO_OUTPUT_FILE" ]]; then rm "eco/$ECO_OUTPUT_FILE"; fi || return
	if [[ -f "eco/$ECO_OUTPUT_FILE2" ]]; then rm "eco/$ECO_OUTPUT_FILE2"; fi || return
	go test -test.fullpath=true -timeout 30s -run ^TestGenBootstrap$ github.com/dylt-dev/dylt/eco -v || return
	go test -test.fullpath=true -timeout 30s -run ^TestGenDecodeDeepTest$ github.com/dylt-dev/dylt/eco -v || return
	go test -test.fullpath=true -timeout 30s -run ^TestDecodeDeepType$ECO_DEPTH github.com/dylt-dev/dylt/eco -v || return
	go fmt "eco/$ECO_OUTPUT_FILE"
	go fmt "eco/$ECO_OUTPUT_FILE2"
}

main
