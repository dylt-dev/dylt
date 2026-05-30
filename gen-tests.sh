#! /usr/bin/env bash

main ()
{
	export ECO_GEN_TESTS=Y
	export ECO_DEPTH=1000
	export ECO_GENNER_FILENAME_PREFIX=deep_genner_
	export ECO_GENNER_TESTNAME_PREFIX=TestDeepGenner
	export ECO_TEST_COUNT=100
	export ECO_TEST_FILENAME_PREFIX=deep_genned_
	export ECO_TEST_TESTNAME_PREFIX=TestDeepGenned
	
	rm eco/$ECO_GENNER_FILENAME_PREFIX*.go
	rm eco/$ECO_TEST_FILENAME_PREFIX*.go
	go test -test.fullpath=true -timeout 300s eco/gen_bootstrap_test.go -v || return

	local tmpGenFolder; tmpGenFolder=$(mktemp --directory --tmpdir -t deepgen) || return
	mv "./eco/$ECO_GENNER_FILENAME_PREFIX"*.go "$tmpGenFolder" || return
	for ((i=0; i<$ECO_TEST_COUNT; i++)); do
		cp "$tmpGenFolder/${ECO_GENNER_FILENAME_PREFIX}${i}_test.go" ./eco || return
		go test -test.fullpath=true -timeout 300s -run "^TestDeepGenner${i}\$" github.com/dylt-dev/dylt/eco -v || return
		go test -test.fullpath=true -timeout 300s -run "^TestDeepGenned${i}\$" "github.com/dylt-dev/dylt/eco" -v || return
		rm "./eco/${ECO_GENNER_FILENAME_PREFIX}${i}_test.go" || return
		rm "./eco/${ECO_TEST_FILENAME_PREFIX}${i}_test.go" || return
	done

	# gofmt -s -w eco
}

main
