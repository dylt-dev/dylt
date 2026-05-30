#! /usr/bin/env bash

main ()
{
	export ECO_GEN_TESTS=1
	
	export ECO_DEPTH=2
	export ECO_GENNER_FILENAME_PREFIX=deep_genner_
	export ECO_GENNER_TESTNAME_PREFIX=TestDeepGenner
	export ECO_TEST_COUNT=1

	# rm "$ECO_GENNER_FILENAME_PREFIX*.go"
	go test -test.fullpath=true -timeout 1s -run ^TestGenBootstrap$ github.com/dylt-dev/dylt/eco -v || return
	# go fmt "eco"
}

main
