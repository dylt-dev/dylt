#! /usr/bin/env bash


main ()
{
	export ECOGEN=Y
    export ECOGEN_BATCHSIZE=1
    export ECOGEN_COUNT=10
    export ECOGEN_DEPTH=1000
    export ECOGEN_FILENAME_PREFIX=Encode
    export ECOGEN_GENFILENAME_PREFIX=EncodeGen
    export ECOGEN_TESTNAME_PREFIX=TestEncode
    export ECOGEN_GENTESTNAME_PREFIX=TestEncodeGen
	
    # Create a tmp folder to stash old tests, and stash em
    dt=$(date +%y.%m.%d)
	local tmpGenFolder; tmpGenFolder=$(mktemp --directory --tmpdir -t "gen.encode.tests.stash.$dt") || return
    mv "eco/$ECOGEN_GENFILENAME_PREFIX"_*_test.go "$tmpGenFolder"
    mv "eco/$ECOGEN_FILENAME_PREFIX"_*_test.go "$tmpGenFolder"

    # Stash testgen files in tmpfolder & gen new ones
    go test -test.fullpath=true -timeout 1800s -run ^TestEncodeBootstrap$ -v -count 1 github.com/dylt-dev/dylt/eco || return

    # Stash test files in tmpfolder & gen new ones
    go test -test.fullpath=true -timeout 1800s -run ^${ECOGEN_GENTESTNAME_PREFIX}_ -v -count 1 github.com/dylt-dev/dylt/eco || return

    # Run genned test files
    go test -test.fullpath=true -timeout 1800s -run ^${ECOGEN_TESTNAME_PREFIX}_ -v -count 1 github.com/dylt-dev/dylt/eco || return
    
}


(return 0 2>/dev/null) || main "$@"
