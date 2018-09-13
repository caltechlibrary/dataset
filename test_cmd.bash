#!/bin/bash

function assert_exists() {
	if [ "$#" != "2" ]; then
		echo "wrong number of parameters for $1, $*"
		exit 1
	fi
	if [[ ! -f "$2" && ! -d "$2" ]]; then
		echo "$1: $2 does not exists"
		exit 1
	fi

}

function assert_equal() {
	if [ "$#" != "3" ]; then
		echo "wrong number of parameters for $1, $*"
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
	if [[ -f "testdata/test1.ds/collection.json" ]]; then
		rm -fR testdata/test1.ds
	fi
	EXT=".exe"
	OS=$(uname)
	if [ "$OS" != "Windows" ]; then
		EXT=""
	fi
	echo "checking for bin/dataset${EXT}"
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
	EXPECTED="OK"
	RESULT=$(bin/dataset init testdata/test1.ds)
	assert_equal "init testdata/test1.ds" "$EXPECTED" "$RESULT"
	assert_exists "collection create" "testdata/test1.ds"
	assert_exists "collection created metadata" "testdata/test1.ds/collection.json"

	export DATASET="testdata/test1.ds"

	# Test create
	EXPECTED="OK"
	RESULT=$(bin/dataset create "${DATASET}" 1 '{"one":1}')
	assert_equal "create 1:" "$EXPECTED" "$RESULT"
	RESULT=$(echo -n '{"two":2}' | bin/dataset create -i - "${DATASET}" 2)
	assert_equal "create 2:" "$EXPECTED" "$RESULT"
	echo '{"three":3}' >"testdata/test3.json"
	RESULT=$(bin/dataset create -i testdata/test3.json "${DATASET}" 3)
	assert_equal "create 3:" "$EXPECTED" "$RESULT"
	echo '{"four":4}' >"testdata/test4.json"
	RESULT=$(bin/dataset create "${DATASET}" 4 testdata/test4.json)
	assert_equal "create 4:" "$EXPECTED" "$RESULT"

	# Test read
	EXPECTED='{"_Key":"1","one":1}'
	RESULT=$(bin/dataset read "${DATASET}" 1)
	assert_equal "read 1:" "$EXPECTED" "$RESULT"
	EXPECTED='{"_Key":"2","two":2}'
	RESULT=$(echo -n '2' | bin/dataset -i - read "${DATASET}")
	assert_equal "read 2:" "$EXPECTED" "$RESULT"

	# Test keys
	EXPECTED="1 2 3 4 "
	RESULT=$(bin/dataset keys "${DATASET}" | sort | tr "\n" " ")
	assert_equal "keys:" "$EXPECTED" "$RESULT"
	EXPECTED="1 "
	RESULT=$(bin/dataset keys "${DATASET}" '(eq .one 1)' | sort | tr "\n" " ")

	if [ -f "testdata/test1.ds/collection.json" ]; then
		rm -fR testdata/test1.ds
	fi
	echo "test_dataset, OK"
}

function test_gsheet() {
    CLIENT_SECRET="${1}"
    echo "test_dataset"
    if [[ ! -f "${CLIENT_SECRET}" ]]; then
        echo "Skipping, could not find ${CLIENT_SECRET}"
        return
    fi
    echo "test_gsheet"
	if [[ -f "etc/test_gsheet.bash" ]]; then
		. "etc/test_gsheet.bash"
	else
		echo "Skipping Google Sheets test, no /etc/test_gsheet.bash found"
		exit 1
	fi
	if [[ ! -s "${CLIENT_SECRET_JSON}" ]]; then
		echo "Skipping test_gsheet(), missing environment varaiable for CLIENT_SECRET_JSON"
		exit 1
	fi
	if [[ "${SPREADSHEET_ID}" == "" ]]; then
		echo "Skipping test_gsheet(), missing environment variable for SPREADSHEET_ID"
		exit 1
	fi
    cd gsheet || exit 1
    go test -client-secret "../${CLIENT_SECRET_JSON}" -spreadsheet-id "${SPREADSHEET_ID}"
    cd ..
	if [[ -d "testdata/test_gsheet.ds" ]]; then
		rm -fR testdata/test_gsheet.ds
	fi
	bin/dataset -nl=false -quiet init "testdata/test_gsheet.ds"
	if [[ "$?" != "0" ]]; then
		echo "Count not initialize testdata/test_gsheet.ds"
		exit 1
	fi
	export DATASET="testdata/test_gsheet.ds"

	bin/dataset -nl=false -quiet create "${DATASET}" "Wilson1930" '{"additional":"Supplemental Files Information:\nGeologic Plate: Supplement 1 from \"The geology of a portion of the Repetto Hills\" (Thesis)\n","description_1":"Supplement 1 in CaltechDATA: Geologic Plate","done":"yes","identifier_1":"https://doi.org/10.22002/D1.638","key":"Wilson1930","resolver":"http://resolver.caltech.edu/CaltechTHESIS:12032009-111148185","subjects":"Repetto Hills, Coyote Pass, sandstones, shales"}'
	if [[ "$?" != "0" ]]; then
		echo "Could not create test record in testdata/test_gsheet.ds"
		exit 1
	fi
	CNT=$(bin/dataset -nl=false count "${DATASET}")
	if [[ "${CNT}" != "1" ]]; then
		echo "Should have one record to export"
		exit 1
	fi

	echo "test_gsheet: export support "
	SHEET_NAME="Sheet1"
	bin/dataset -nl=false -quiet export -client-secret "${CLIENT_SECRET_JSON}" "${DATASET}" "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' \
		'.done,.key,.resolver,.subjects,.additional,.identifier_1,.description_1' \
		'Done,Key,Resolver,Subjects,Additional,Identifier 1,Description 1'
	if [[ "$?" != "0" ]]; then
		echo "Count not export-gsheet"
		exit 1
	fi

	echo "test_gsheet: import support "
	bin/dataset -nl=false -quiet import -client-secret "${CLIENT_SECRET_JSON}" "${DATASET}" "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' 2
	if [[ "$?" == "0" ]]; then
		echo "Should NOT be able to import-gsheet over our existing collection without -overwrite"
		exit 1
	fi

	bin/dataset -nl=false -quiet import -overwrite -client-secret "${CLIENT_SECRET_JSON}" "${DATASET}" "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' 2
	if [[ "$?" != "0" ]]; then
		echo "Should be able to import-gsheet over our existing collection with -overwrite"
		exit 1
	fi

    # Check to see if this throws error correctly, i.e. should have exit code 1
	SHEET_NAME="Sheet2"
	bin/dataset -nl=false -quiet export -client-secret "${CLIENT_SECRET_JSON}" "${DATASET}" "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' true \
		'true,.done,.key,.QT_resolver,.subjects,.additional[],.identifier_1,.description_1' \
		'Done,Key,Resolver,Subjects,Additional,Identifier 1,Description 1'
	if [[ "$?" != "1" ]]; then
		echo "Count should throw error for bad dotpath in export-gsheet"
		exit 1
	fi

	echo "test_gsheet, OK"
}


function test_issue19() {
	echo "test_issue19"
	if [[ -d "testdata/test_issue19.ds" ]]; then
		rm -fR testdata/test_issue19.ds
	fi
	bin/dataset -nl=false -quiet init "testdata/test_issue19.ds"
	bin/dataset -nl=false -quiet create testdata/test_issue19.ds freda '{"name":"freda","email":"freda@inverness.example.org","try":1}'
	if [[ "$?" != "0" ]]; then
		echo "Failed, should be able to create the record in an empty collection"
		exit 1
	fi

	# Now try creating the record again without -overwrite
    echo "NOTE: expecting an error message on next line"
	bin/dataset -nl=false -quiet create testdata/test_issue19.ds freda '{"name":"freda","email":"freda@inverness.example.org","try":2}'
	if [[ "$?" != "1" ]]; then
        echo "Expected return value 1, got $?"
		echo "Failed, should NOT be able to create the record when it exists in a collection without -overwrite"
		exit 1
	fi

	# Now try to create the record with -overwrite
	bin/dataset -nl=false -quiet create -overwrite testdata/test_issue19.ds freda '{"name":"freda","email":"freda@inverness.example.org","try":3}'
	if [[ "$?" != "0" ]]; then
		echo "Failed, should be able to create the record with -overwrite!"
		exit 1
	fi

	echo "test_issue19, OK"
	rm -fR "testdata/test_issue19.ds"
}

function test_readme () {
    echo "test_readme"
    mkdir -p testdata
    if [[ -d "testdata/mystuff.ds" ]]; then
        rm -fR testdata/mystuff.ds
    fi

    # Create a collection "testdata/mystuff.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    bin/dataset -quiet -nl=false init testdata/mystuff.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (206): could not init mystuff.ds'
        exit 1
    fi
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    bin/dataset -quiet -nl=false create testdata/mystuff.ds freda '{"name":"freda","email":"freda@inverness.example.org"}'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (216): could not create freda.json'
        exit 1
    fi
    # If successful then you should see an OK otherwise an error message

    # Make sure we have a record called freda
    bin/dataset -quiet -nl="false" haskey testdata/mystuff.ds freda > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (222): (failed) testdata/mystuff.ds haskey freda'
        exit 1
    fi


    # Read a JSON document
    bin/dataset -quiet -nl=false read testdata/mystuff.ds freda > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (232): could not read freda.json'
        exit 1
    fi
    
    # Path to JSON document
    bin/dataset -quiet -nl=false path testdata/mystuff.ds freda > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (237): could not path freda.json'
        exit 1
    fi

    # Update a JSON document
    bin/dataset -quiet -nl=false update testdata/mystuff.ds freda '{"name":"freda","email":"freda@zbs.example.org", "count": 2}'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (244): could not update freda.json'
        exit 1
    fi
    
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    bin/dataset -quiet -nl=false keys testdata/mystuff.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (253): could not keys'
        exit 1
    fi

    # Get keys filtered for the name "freda"
    bin/dataset -nl=false -quiet keys testdata/mystuff.ds '(eq .name "freda")' > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (260): could not keys'
        exit 1
    fi

    # Join freda-profile.json with "freda" adding unique key/value pairs
    cat << EOT > testdata/freda-profile.json
{"name": "little freda", "office": "SFL", "count": 3}
EOT

    bin/dataset -quiet -nl=false join testdata/mystuff.ds freda testdata/freda-profile.json
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (271): could not join update'
        exit 1
    fi

    # Join freda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from freda-profile.json
    cat << EOT > testdata/freda-profile.json
{"name": "little freda", "office": "SFL", "count": 4}
EOT

    bin/dataset -quiet -nl=false join -overwrite testdata/mystuff.ds freda testdata/freda-profile.json
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (283): could not join overwrite'
        exit 1
    fi


    # Delete a JSON document
    bin/dataset -quiet -nl=false delete testdata/mystuff.ds freda
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (290): could not join overwrite'
        exit 1
    fi

    # Import from a CSV file
    cat << EOT > testdata/my-data.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
EOT

    bin/dataset -quiet -nl=false import testdata/mystuff.ds testdata/my-data.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_readme (302): (failed) import testdata/mystuff.ds testdata/my-data.csv 1'
        exit 1
    fi

    echo "test_readme, OK"
    # To remove the collection just use the Unix shell command
    rm -fR testdata/mystuff.ds
    rm testdata/freda-profile.json
    rm testdata/my-data.csv
}

function test_getting_started() {
    echo "test_getting_started"
    if [[ -d "testdata/FavoriteThings.ds" ]]; then
        rm -fR testdata/FavoriteThings.ds
    fi
    bin/dataset -quiet -nl=false init testdata/FavoriteThings.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not init testdata/FavoriteThings.ds'
        exit 1
    fi

    bin/dataset -quiet -nl=false create testdata/FavoriteThings.ds beverage '{"thing":"coffee"}'
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not testdata/FavoriteThings.ds create beverage'
        exit 1
    fi

    bin/dataset -quiet -nl=false read testdata/FavoriteThings.ds beverage > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not testdata/FavoriteThings.ds read beverage'
        exit 1
    fi

    bin/dataset -quiet -nl=false keys testdata/FavoriteThings.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not testdata/FavoriteThings.ds keys'
        exit 1
    fi

    cat << EOT > testdata/jazz-notes.json
{
    "songs": ["Blue Rondo al la Turk", "Bernie's Tune", "Perdido"],
    "pianist": [ "Dave Brubeck" ],
    "trumpet": [ "Dirk Fischer", "Dizzy Gillespie" ]
}
EOT
    bin/dataset -quiet -nl=false create testdata/FavoriteThings.ds "jazz-notes" testdata/jazz-notes.json
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not create jazz-notes'
        exit 1
    fi

    bin/dataset -quiet -nl=false keys testdata/FavoriteThings.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not keys'
        exit 1
    fi

    bin/dataset -quiet -nl=false read testdata/FavoriteThings.ds beverage jazz-notes > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not read multiple keys'
        exit 1
    fi

    # Cleanup after tests
    rm -fR testdata/FavoriteThings.ds
    rm testdata/jazz-notes.json
    echo "test_getting_started, OK"
}

function test_attachments() {
    echo 'test_attachments'
    if [[ -d "testdata/mydata.ds" ]]; then
        rm -fR testdata/mydata.ds
    fi
    bin/dataset -quiet -nl=false init testdata/mydata.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (375): could not testdata/mydata.ds init'
        exit 1
    fi

    cat << EOT > testdata/freda.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
EOT

    cat << EOT > testdata/mojo.csv
Name,EMail,Office,Count
mojo,mojo.sam@sams-splace.example.org,piano,2
EOT

    bin/dataset -quiet -nl=false import testdata/mydata.ds testdata/freda.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (389): (failed) testdata/mydata.ds import-csv testdata/freda.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false attach testdata/mydata.ds freda testdata/freda.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (395): (failed) testdata/mydata.ds attach freda testdata/freda.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/mydata.ds testdata/mojo.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (399): (failed) testdata/mydata.ds import-csv testdata/mojo.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false attach testdata/mydata.ds mojo testdata/mojo.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (404): (failed) testdata/mydata.ds attach testdata/mojo.csv'
        exit 1
    fi
    bin/dataset -quiet -nl=false attachments testdata/mydata.ds mojo > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (410): (failed) testdata/mydata.ds attachments mojo'
        exit 1
    fi
    if [[ -f "testdata/mojo.csv" ]]; then
        rm testdata/mojo.csv
    fi
    bin/dataset -quiet -nl=false detach testdata/mydata.ds mojo testdata/mojo.csv > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (417): (failed) testdata/mydata.ds attachments mojo testdata/mojo.csv'
        exit 1
    fi
    if [[ ! -f "testdata/mojo.csv" ]]; then
        echo 'test_attachments (417): (failed) testdata/mydata.ds detach mojo testdata/mojo.csv'
        exit 1
    fi
    bin/dataset -quiet -nl=false prune testdata/mydata.ds freda testdata/freda.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments (426): (failed) testdata/mydata.ds prune freda testdata/freda.csv'
        exit 1
    fi

    # Success, cleanup our test data
    rm testdata/freda.csv testdata/mojo.csv
    rm -fR testdata/mydata.ds
	echo "test_attachments, OK"
}

function test_check_and_repair() {
    echo 'test_check_and_repair'
    if [[ -d "testdata/myfix.ds" ]]; then
        rm -fR testdata/myfix.ds
    fi
    cat << EOT > testdata/myfix.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
mojo,mojo.sam@sams-splace.example.org,piano,2
EOT
    bin/dataset -quiet -nl=false init testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) init testdata/myfix.ds'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/myfix.ds testdata/myfix.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) import testdata/myfix.ds testdata/myfix.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false check testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) check testdata/myfix.ds'
        exit 1
    fi
    
    CNT=$(bin/dataset count testdata/myfix.ds)
    echo "NOTE: expecting ${CNT} error messages on following lines"
    echo '{}' > testdata/myfix.ds/collection.json
    bin/dataset -quiet -nl=false check testdata/myfix.ds
    if [[ "$?" == "0" ]]; then
        echo 'test_check_and_repair: (failed, expected exit code 1) testdata/myfix.ds check'
        exit 1
    fi
    echo "NOTE: Initiating a repair"
    bin/dataset -quiet -nl=false repair testdata/myfix.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) repair testdata/myfix.ds'
        exit 1
    fi
    echo "NOTE: Expecting OK on next line if repair worked"
    bin/dataset check testdata/myfix.ds 
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) check testdata/myfix.ds'
        exit 1
    fi

    echo "NOTE: Testing migration from pairtree to buckets, expect OK"
    bin/dataset migrate testdata/myfix.ds buckets
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) migrate testdata/myfix.ds buckets'
        exit 1
    fi

    echo "NOTE: Testing migration from buckets to pairtree, expect OK"
    bin/dataset migrate testdata/myfix.ds pairtree
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) migrate testdata/myfix.ds pairtree'
        exit 1
    fi

    echo "NOTE: Final Check, expecting OK"
    bin/dataset check testdata/myfix.ds 
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) check testdata/myfix.ds'
        exit 1
    fi

   
    # Success, cleanup
    rm -fR testdata/myfix.ds
    rm testdata/myfix.csv
    echo 'test_check_and_repair, OK'
}

function test_count() {
    echo 'test_count'
    if [[ -d "testdata/count.ds" ]]; then
        rm -fR testdata/count.ds
    fi
    cat << EOT > testdata/count.csv
Name,EMail,Office,Count,published
freda,freda@inverness.example.edu,4th Tower,1,true
mojo,mojo.sam@sams-splace.example.org,piano,2,false
EOT
 
    if [[ ! -f "testdata/count.csv" ]]; then
        echo 'test_count: (failed) could not create testdata/count.csv'
        exit 1
    fi

    bin/dataset -quiet -nl=false init testdata/count.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) init testdata/count.ds'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/count.ds testdata/count.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) import testdata/count.ds testdata/count.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false count testdata/count.ds > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) testdata/count.ds count'
        exit 1
    fi
    bin/dataset -quiet -nl=false count testdata/count.ds '(eq .published true)' > /dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) count testdata/count.ds "(eq .published true)"'
        exit 1
    fi

    # Success, cleanup
    rm -fR testdata/count.ds
    rm testdata/count.csv
    echo 'test_count, OK'
}

function test_search() {
    echo 'test_search'
    if [[ -d "testdata/search.ds" ]]; then
        rm -fR testdata/search.ds
    fi
    cat << EOT > testdata/search.csv
id,title,type,published,author
a1,4th Tower of Inverness,audio play,true,Tom Lopez
n1,"20,000 leagues under the Sea",novel,true,Jules Verne
s1,Our Person in Avalon,screenplay,false,R. S. Doiel
EOT
    
    cat << EOT > testdata/search.json
{
    "title": {
        "object_path": ".title"
    },
    "type": {
        "object_path": ".type"
    },
    "published": {
        "object_path": ".published"
    },
    "author": {
        "object_path": ".author"
    }

}
EOT

    bin/dataset -quiet -nl=false init testdata/search.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) init testdata/search.ds'
        exit 1
    fi
    bin/dataset -quiet -nl=false import testdata/search.ds testdata/search.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) testdata/search.ds import-csv testdata/search.csv 1'
        exit 1
    fi
    # FIXME: this needs to use a frame to define search
    bin/dataset -quiet -nl=false indexer testdata/search.ds testdata/search.json testdata/search.bleve
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) testdata/search.ds indexer testdata/search.json testdata/search.bleve'
        exit 1
    fi
    bin/dataset -quiet -nl=false find testdata/search.bleve 'screenplay' 
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) testdata/search.ds find "screenplay"'
        exit 1
    fi
    bin/dataset -quiet -nl=false find testdata/search.bleve '+published:true' 
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) testdata/search.ds find "+published:true"'
        exit 1
    fi
    bin/dataset -quiet -nl=false find testdata/search.bleve 'Tower'
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) testdata/search.ds find "Tower"'
        exit 1
    fi
    echo 'a1' > testdata/list.keys
    bin/dataset -quiet -nl=false deindexer testdata/search.bleve testdata/list.keys
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) deindexer testdata/search.bleve testdata/list.keys' 
        exit 1
    fi

    # Success, cleanup
    rm -fR testdata/search.ds
    rm -fR testdata/search.bleve
    rm testdata/search.json
    rm testdata/search.csv
    rm testdata/list.keys
    echo 'test_search, OK'
}

function test_import_export() {
    echo 'test_import_export'
    if [[ -d "testdata/pubs.ds" ]]; then
        rm -fR "testdata/pubs.ds"
    fi
    cat << EOT > testdata/in.csv
id,title,type,date_type,date
44088,Application of a laser induced fluorescence model to the numerical simulation of detonation waves in hydrogen-oxygen-diluent mixtures,article,published,2014-04-04
46001,Leaderless Deterministic Chemical Reaction Networks,book_section,published,2013
61958,"Observation of the 14 MeV resonance in ^(12)C(p, p)^(12)C with molecular ion beams",article,published,1972-02-21
21682,Effect of scale size on a rocket engine with suddenly frozen nozzle flow,article,published,1961-03
39470,A Review of the Dynamics of Cavitating Pumps,article,published,2012-11-26
62459,An Experimental Investigation of the Flow over Blunt-Nosed Cones at a Mach Number of 5.8,monograph,completed,1956-06-15
80289,CIT-4: The first synthetic analogue of brewsterite,article,published,1997-08
80630,The Allocation of a Shared Resource Within an Organization,monograph,completed,1995-01
8488,Non-Gaussian covariance of CMB B modes of polarization and parameter degradation,article,published,2007-04-15
EOT

    bin/dataset -quiet -nl=false init testdata/pubs.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) init testdata/pubs.ds'
        exit 1
    fi

    bin/dataset -quiet -nl=false import testdata/pubs.ds testdata/in.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) testdata/pubs.ds import-csv testdata/in.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false keys testdata/pubs.ds >/dev/null
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) testdata/pubs.ds keys'
        exit 1
    fi
    #FIXME: export uses a frame to define exported content
    bin/dataset -quiet -nl=false frame -a testdata/pubs.ds outframe \
         "._Key" ".title" ".type" ".date_type" ".date" > /dev/null
    bin/dataset -quiet -nl=false frame-labels  testdata/pubs.ds outframe \
        'EPrint ID' 'Title' 'Type' 'Date Type' 'Date'
    bin/dataset -quiet -nl=false export testdata/pubs.ds outframe "testdata/out.csv"
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) export testdata/pubs.ds outframe testdata/out.csv'
        exit 1
    fi

    # Success, cleanup
    rm -fR testdata/data.ds
    rm testdata/in.csv
    rm testdata/out.csv
    echo 'test_import_export, OK'
}

echo "Testing command line tools"
test_dataset
test_issue19
test_readme
test_getting_started
test_attachments
test_count
test_import_export
#NOTE: test will be skip if there is no etc/client_secret.json found
test_gsheet etc/client_secret.json
test_search
test_check_and_repair
echo 'PASS'
echo "Ok $(basename "$0")"
