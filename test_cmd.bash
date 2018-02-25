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
    if [ -f "test1.ds/collection.json" ]; then
        rm -fR test1.ds
    fi
    EXT=".exe"
    OS=$(uname)
    if [ "$OS" != "Windows" ]; then
        EXT=""
    fi
    echo "Testing for bin/dataset${EXT}"
    if [[ ! -f "bin/dataset${EXT}" || ! -f "cmd/dataset/assets.go" ]]; then
        # We need to build
	    pkgassets -o cmd/dataset/assets.go \
            -p main -ext=".md" -strip-prefix="/" \
            -strip-suffix=".md" \
            Examples examples/dataset \
            Help docs/dataset
        go build -o "bin/dataset${EXT}" cmd/dataset/dataset.go cmd/dataset/assets.go
    fi

    # Test init
    EXPECTED='export DATASET=test1.ds'
    RESULT=$(bin/dataset init test1.ds)
    assert_equal "init test1.ds" "$EXPECTED" "$RESULT"
    assert_exists "collection create" "test1.ds"
    assert_exists "collection created metadata" "test1.ds/collection.json"
    export DATASET="test1.ds"

    # Test create 
    EXPECTED="OK"
    RESULT=$(bin/dataset create 1 '{"one":1}')
    assert_equal "create 1:" "$EXPECTED" "$RESULT" 
    RESULT=$(echo -n '{"two":2}' | bin/dataset -i - create 2)
    assert_equal "create 2:" "$EXPECTED" "$RESULT" 
    echo '{"three":3}' >"testdata/test3.json"
    RESULT=$(bin/dataset -i testdata/test3.json create 3)
    assert_equal "create 3:" "$EXPECTED" "$RESULT" 

    # Test read
    EXPECTED='{"_Key":"1","one":1}'
    RESULT=$(bin/dataset read 1)
    assert_equal "read 1:" "$EXPECTED" "$RESULT"
    EXPECTED='{"_Key":"2","two":2}'
    RESULT=$(echo -n '2' | bin/dataset -i - read)
    assert_equal "read 1:" "$EXPECTED" "$RESULT"

    # Test keys
    EXPECTED="1 2 3 "
    RESULT=$(bin/dataset keys | sort | tr "\n" " ")
    assert_equal "keys:" "$EXPECTED" "$RESULT"
    EXPECTED="1 "
    RESULT=$(bin/dataset keys '(eq .one 1)' | sort | tr "\n" " ")

    if [ -f "test1.ds/collection.json" ]; then
        rm -fR test1.ds
    fi
    echo "Test dataset successful"
}

function test_gsheet() {
    if [[ -f "etc/test_gsheet.bash" ]]; then
        . "etc/test_gsheet.bash"
    else
        echo "Skipping Google Sheets test, no /etc/test_gsheet.bash found"
        exit 1
    fi
    if [[ ! -s "${CLIENT_SECRET_JSON}" ]]; then
        echo "Missing environment varaiable for CLIENT_SECRET_JSON"
        exit 1
    fi
    if [[ "${SPREADSHEET_ID}" = "" ]]; then
        echo "Missing environment variable for SPREADSHEET_ID"
        exit 1
    fi
    echo "Testing Google Sheets support"
    if [[ -d "test_gsheet.ds" ]]; then
        rm -fR test_gsheet.ds
    fi
    bin/dataset init "test_gsheet.ds"
    if [[ "$?" != "0" ]]; then
        echo "Count not initialize test_gsheet.ds"
        exit 1
    fi
    export DATASET="test_gsheet.ds"

    bin/dataset create test '{"additional":"Supplemental Files Information:\nGeologic Plate: Supplement 1 from \"The geology of a portion of the Repetto Hills\" (Thesis)\n","description_1":"Supplement 1 in CaltechDATA: Geologic Plate","done":"yes","identifier_1":"https://doi.org/10.22002/D1.638","key":"Wilson1930","resolver":"http://resolver.caltech.edu/CaltechTHESIS:12032009-111148185","subjects":"Repetto Hills, Coyote Pass, sandstones, shales"}'
    if [[ "$?" != "0" ]]; then
        echo "Could not create test record in test_gsheet.ds"
        exit 1
    fi
    CNT=$(bin/dataset count)
    if [[ "${CNT}" != "1" ]]; then
        echo "Should have one record to export"
        exit 1
    fi

    bin/dataset -client-secret "${CLIENT_SECRET_JSON}" export-gsheet "${SPREADSHEET_ID}" 'Sheet1' 'A1:CZ' true \
        '.done,.key,.resolver,.subjects,.additional,.identifier_1,.description_1' \
        'Done,Key,Resolver,Subjects,Additional,Identifier 1,Description 1'
    if [[ "$?" != "0" ]]; then
        echo "Count not export-gsheet"
        exit 1
    fi

    echo "Test dataset gsheet support successful"
}

function test_issue15() {
    if [[ -d "test_issue15.ds" ]]; then
        rm -fR test_issue15.ds
    fi
    bin/dataset init "test_issue15.ds"
    bin/dataset -c test_issue15.ds create freda '{"name":"freda","email":"freda@inverness.example.org"}'
    I=$(bin/dataset -c test_issue15.ds count)
    if [[ "$I" != "1" ]]; then
        echo "Failed to add freda record test_issue15.ds"
        exit 1
    fi
    K=$(bin/dataset -nl=false -c test_issue15.ds keys '(eq "freda" .name)')
    if [[ "$K" != "freda" ]]; then
        echo "Should have one key, freda, in test_issue15.ds"
        exit 1
    fi
    V=$(bin/dataset -nl=false -c test_issue15.ds extract 'true' '.name')
    if [[ "$V" != "freda" ]]; then
        echo "Should extract one name, freda, in test_issue15.ds $V"
        exit 1
    fi
    echo "Test issue 15 fix OK"
    rm -fR "test_issue15.ds"
}

echo "Testing command line tools"
test_dataset
test_gsheet
test_issue15
echo 'Success!'
