#!/bin/bash


function assert_exists() {
    if [ "$#" != "2" ]; then
        echo "wrong number of parameters for $1, $@"
        exit 1
    fi
    if [[ ! -f "$2" && ! -d "$2" ]]; then
        echo "$1: $2 does not exists"
        exit 1
    fi

}

function assert_equal() {
    if [ "$#" != "3" ]; then
        echo "wrong number of parameters for $1, $@"
        exit 1
    fi
    if [ "$2" != "$3" ]; then
        echo "$1: expected |$2| got |$3|"
        exit 1
    fi
}



#
# Tests
#

function test_dataset() {
    if [ -f "test1/collection.json" ]; then
        rm -fR test1
    fi
    EXT=".exe"
    OS=$(uname)
    if [ "$OS" != "Windows" ]; then
        EXT=""
    fi
    echo "Testing for bin/dataset${EXT}"
    if [[ ! -f "bin/dataset${EXT}" || ! -f "cmds/dataset/assets.go" ]]; then
        # We need to build
	    pkgassets -o cmds/dataset/assets.go \
            -p main -ext=".md" -strip-prefix="/" \
            -strip-suffix=".md" \
            Examples examples/dataset \
            Help docs/dataset
        go build -o "bin/dataset${EXT}" cmds/dataset/dataset.go cmds/dataset/assets.go
    fi

    # Test init
    EXPECTED='export DATASET=test1'
    RESULT=$(dataset init test1)
    assert_equal "init test1" "$EXPECTED" "$RESULT"
    assert_exists "collection create" "test1"
    assert_exists "collection created metadata" "test1/collection.json"
    export DATASET="test1"

    # Test create 
    EXPECTED="OK"
    RESULT=$(dataset create 1 '{"one":1}')
    assert_equal "create 1:" "$EXPECTED" "$RESULT" 
    RESULT=$(echo -n '{"two":2}' | dataset -i - create 2)
    assert_equal "create 2:" "$EXPECTED" "$RESULT" 

    # Test read
    EXPECTED='{"one":1}'
    RESULT=$(dataset read 1)
    assert_equal "read 1:" "$EXPECTED" "$RESULT"
    EXPECTED='{"two":2}'
    RESULT=$(echo -n '2' | dataset -i - read)
    assert_equal "read 1:" "$EXPECTED" "$RESULT"

    # Test keys
    EXPECTED="1 2 "
    RESULT=$(dataset keys | sort | tr "\n" " ")
    assert_equal "keys:" "$EXPECTED" "$RESULT"
    EXPECTED="1 "
    RESULT=$(dataset keys '(eq .one 1)' | sort | tr "\n" " ")

    if [ -f "test1/collection.json" ]; then
        rm -fR test1
    fi
    echo "Test dataset successful"
}

test_dataset
echo 'Success!'
