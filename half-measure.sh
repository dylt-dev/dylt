#! /usr/bin/env bash

main ()
{
	export ECO_GEN_TESTS=1
	export ECO_TEST_FILENAME_PREFIX=deep_genned_
	export ECO_TEST_TESTNAME_PREFIX=TestDeepGenned
	# export ECO_DEPTH=1000
	export ECO_GENNED_TEST_FOLDER=./biguns
	# export ECO_TEST_COUNT=1

	# if [[ -f "eco/$ECO_STAGE1_FILE" ]]; then rm "eco/$ECO_STAGE1_FILE"; fi || return
	mkdir -p "eco/$ECO_GENNED_TEST_FOLDER" || return
	# go test -test.fullpath=true -timeout 30s -run ^TestGenBootstrap$ github.com/dylt-dev/dylt/eco -v || return
	go test -test.fullpath=true -timeout 30s -run ^TestDeepGenner github.com/dylt-dev/dylt/eco -v || return
	# go test -test.fullpath=true -timeout 30s -run ^TestDecodeDeepType$ECO_DEPTH github.com/dylt-dev/dylt/eco -v || return
	# go fmt "eco/$ECO_STAGE1_FILE"
	# go fmt "eco/$ECO_GENNED_TEST_FOLDER"
}

main
