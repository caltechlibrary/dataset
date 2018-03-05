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
	if [[ -f "test1.ds/collection.json" ]]; then
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
	EXPECTED="OK"
	RESULT=$(bin/dataset init test1.ds)
	assert_equal "init test1.ds" "$EXPECTED" "$RESULT"
	assert_exists "collection create" "test1.ds"
	assert_exists "collection created metadata" "test1.ds/collection.json"

	# Set environment and then continue with tests
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
	echo '{"four":4}' >"testdata/test4.json"
	RESULT=$(bin/dataset create 4 testdata/test4.json)
	assert_equal "create 4:" "$EXPECTED" "$RESULT"

	# Test read
	EXPECTED='{"_Key":"1","one":1}'
	RESULT=$(bin/dataset read 1)
	assert_equal "read 1:" "$EXPECTED" "$RESULT"
	EXPECTED='{"_Key":"2","two":2}'
	RESULT=$(echo -n '2' | bin/dataset -i - read)
	assert_equal "read 1:" "$EXPECTED" "$RESULT"

	# Test keys
	EXPECTED="1 2 3 4 "
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
	if [[ "${SPREADSHEET_ID}" == "" ]]; then
		echo "Missing environment variable for SPREADSHEET_ID"
		exit 1
	fi
	if [[ -d "test_gsheet.ds" ]]; then
		rm -fR test_gsheet.ds
	fi
	bin/dataset -nl=false -quiet init "test_gsheet.ds"
	if [[ "$?" != "0" ]]; then
		echo "Count not initialize test_gsheet.ds"
		exit 1
	fi
	export DATASET="test_gsheet.ds"

	bin/dataset -nl=false -quiet create "Wilson1930" '{"additional":"Supplemental Files Information:\nGeologic Plate: Supplement 1 from \"The geology of a portion of the Repetto Hills\" (Thesis)\n","description_1":"Supplement 1 in CaltechDATA: Geologic Plate","done":"yes","identifier_1":"https://doi.org/10.22002/D1.638","key":"Wilson1930","resolver":"http://resolver.caltech.edu/CaltechTHESIS:12032009-111148185","subjects":"Repetto Hills, Coyote Pass, sandstones, shales"}'
	if [[ "$?" != "0" ]]; then
		echo "Could not create test record in test_gsheet.ds"
		exit 1
	fi
	CNT=$(bin/dataset -nl=false count)
	if [[ "${CNT}" != "1" ]]; then
		echo "Should have one record to export"
		exit 1
	fi

	echo "Test gsheet export support "
	SHEET_NAME="Sheet1"
	bin/dataset -nl=false -quiet -client-secret "${CLIENT_SECRET_JSON}" export-gsheet "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' true \
		'.done,.key,.resolver,.subjects,.additional,.identifier_1,.description_1' \
		'Done,Key,Resolver,Subjects,Additional,Identifier 1,Description 1'
	if [[ "$?" != "0" ]]; then
		echo "Count not export-gsheet"
		exit 1
	fi

	echo "Test import-gsheet support "
	bin/dataset -nl=false -quiet -client-secret "${CLIENT_SECRET_JSON}" import-gsheet "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' 2
	if [[ "$?" == "0" ]]; then
		echo "Should NOT be able to import-gsheet over our existing collection without -overwrite"
		exit 1
	fi

	bin/dataset -nl=false -quiet -overwrite -client-secret "${CLIENT_SECRET_JSON}" import-gsheet "${SPREADSHEET_ID}" "${SHEET_NAME}" 'A1:CZ' 2
	if [[ "$?" != "0" ]]; then
		echo "Should be able to import-gsheet over our existing collection with -overwrite"
		exit 1
	fi

	echo "Test gsheet support successful"
}

function test_issue15() {
    #unset DATASET
	CNAME="test_issue15.ds"
	if [[ -d "$CNAME" ]]; then
		rm -fR $CNAME
	fi
	bin/dataset -nl=false -quiet init $CNAME
	if [[ ! -d "$CNAME" ]]; then
		echo "Should have created $CNAME"
		exit 1
	fi
	bin/dataset -quiet -nl=false "$CNAME" create freda '{"name":"freda","email":"freda@inverness.example.org"}'
	I="$(bin/dataset -nl=false "$CNAME" count)"
	if [[ "$I" != "1" ]]; then
		echo "Failed to add freda record $CNAME"
		exit 1
	fi
	K="$(bin/dataset -nl=false "$CNAME" keys '(eq "freda" .name)')"
	if [[ "$K" != "freda" ]]; then
		echo "Should have one key, freda, in $CNAME"
		exit 1
	fi
	V="$(bin/dataset -nl=false "$CNAME" extract 'true' '.name')"
	if [[ "$V" != "freda" ]]; then
		echo "Should extract one name, freda, in $CNAME $V"
		exit 1
	fi
	echo "Test issue 15 fix OK"
	rm -fR "$CNAME"
}

function test_issue19() {
	if [[ -d "test_issue19.ds" ]]; then
		rm -fR test_issue19.ds
	fi
	bin/dataset -nl=false -quiet init "test_issue19.ds"
	bin/dataset -nl=false -quiet -c test_issue19.ds create freda '{"name":"freda","email":"freda@inverness.example.org","try":1}'
	if [[ "$?" != "0" ]]; then
		echo "Failed, should be able to create the record in an empty collection"
		exit 1
	fi

	# Now try creating the record again without -overwrite
	bin/dataset -nl=false -quiet -c test_issue19.ds create freda '{"name":"freda","email":"freda@inverness.example.org","try":2}'
	if [[ "$?" == "0" ]]; then
		echo "Failed, should NOT be able to create the record when it exists in an empty collection without -overwrite"
		exit 1
	fi

	# Now try to create the record with -overwrite
	bin/dataset -nl=false -quiet -overwrite -c test_issue19.ds create freda '{"name":"freda","email":"freda@inverness.example.org","try":2}'
	if [[ "$?" != "0" ]]; then
		echo "Failed, should be able to create the record with -overwite!"
		exit 1
	fi

	echo "Test issue 19 fix OK"
	rm -fR "test_issue19.ds"
}

function test_readme () {
    echo "Tests from README.md"
    if [[ -d "mystuff.ds" ]]; then
        rm -fR mystuff.ds
    fi

    # Create a collection "mystuff.ds", the ".ds" lets the bin/dataset command know that's the collection to use. 
    bin/dataset -quiet -nl=false mystuff.ds init
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not init mystuff.ds'
        exit 1
    fi
    # if successful then you should see an OK otherwise an error message

    # Create a JSON document 
    bin/dataset -quiet -nl=false mystuff.ds create freda '{"name":"freda","email":"freda@inverness.example.org"}'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not create freda.json'
        exit 1
    fi
    # If successful then you should see an OK otherwise an error message

    # Make sure we have a record called freda
    bin/dataset -quiet -nl="false" mystuff.ds haskey freda
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: (failed) mystuff.ds haskey freda'
        exit 1
    fi


    # Read a JSON document
    bin/dataset -quiet -nl=false mystuff.ds read freda
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not read freda.json'
        exit 1
    fi
    
    # Path to JSON document
    bin/dataset -quiet -nl=false mystuff.ds path freda
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not path freda.json'
        exit 1
    fi

    # Update a JSON document
    bin/dataset -quiet -nl=false mystuff.ds update freda '{"name":"freda","email":"freda@zbs.example.org", "count": 2}'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not update freda.json'
        exit 1
    fi
    
    # If successful then you should see an OK or an error message

    # List the keys in the collection
    bin/dataset -quiet -nl=false mystuff.ds keys
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not keys'
        exit 1
    fi

    # Get keys filtered for the name "freda"
    bin/dataset -nl=false -quiet mystuff.ds keys '(eq .name "freda")'
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not keys'
        exit 1
    fi

    # Join freda-profile.json with "freda" adding unique key/value pairs
    cat << EOT > freda-profile.json
{"name": "little freda", "office": "SFL", "count": 3}
EOT

    bin/dataset -quiet -nl=false mystuff.ds join append freda freda-profile.json
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not join update'
        exit 1
    fi

    # Join freda-profile.json overwriting in commont key/values adding unique key/value pairs
    # from freda-profile.json
    cat << EOT > freda-profile.json
{"name": "little freda", "office": "SFL", "count": 4}
EOT

    bin/dataset -quiet -nl=false mystuff.ds join overwrite freda freda-profile.json
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not join overwrite'
        exit 1
    fi

    # Delete a JSON document
    bin/dataset -quiet -nl=false mystuff.ds delete freda
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: could not join overwrite'
        exit 1
    fi

    # Import from a CSV file
    cat << EOT > my-data.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
EOT

    bin/dataset -quiet -nl=false mystuff.ds "import-csv" my-data.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_readme: (failed) mystuff.ds import-csv my-data.csv 1'
        exit 1
    fi

    # To remove the collection just use the Unix shell command
    rm -fR mystuff.ds
    rm freda-profile.json
    rm my-data.csv
}

function test_getting_started() {
    echo "Tests from Getting Started with Dataset"
    if [[ -d "FavoriteThings.ds" ]]; then
        rm -fR FavoriteThings.ds
    fi
    bin/dataset -quiet -nl=false init FavoriteThings.ds
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not init FavoriteThings.ds'
        exit 1
    fi

    bin/dataset -quiet -nl=false FavoriteThings.ds create beverage '{"thing":"coffee"}'
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not FavoriteThings.ds create beverage'
        exit 1
    fi

    bin/dataset -quiet -nl=false FavoriteThings.ds read beverage
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not FavoriteThings.ds read beverage'
        exit 1
    fi

    bin/dataset -quiet -nl=false FavoriteThings.ds keys
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not FavoriteThings.ds keys'
        exit 1
    fi

    cat << EOT > jazz-notes.json
{
    "songs": ["Blue Rondo al la Turk", "Bernie's Tune", "Perdido"],
    "pianist": [ "Dave Brubeck" ],
    "trumpet": [ "Dirk Fischer", "Dizzy Gillespie" ]
}
EOT
    bin/dataset -quiet -nl=false FavoriteThings.ds create "jazz-notes" jazz-notes.json
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not create jazz-notes'
        exit 1
    fi

    bin/dataset -quiet -nl=false FavoriteThings.ds keys
    if [[ "$?" != "0" ]]; then
        echo 'test_getting_started: could not keys'
        exit 1
    fi

    bin/dataset -quiet -nl=false FavoriteThings.ds list beverage jazz-notes
    # Cleanup after tests
    rm -fR FavoriteThings.ds
    rm jazz-notes.json
}

function test_attachments() {
    echo 'Test attachments'
    if [[ -d "mydata.ds" ]]; then
        rm -fR mydata.ds
    fi
    bin/dataset -quiet -nl=false mydata.ds init
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: could not mydata.ds init'
        exit 1
    fi

    cat << EOT > freda.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
EOT

    cat << EOT > mojo.csv
Name,EMail,Office,Count
mojo,mojo.sam@sams-splace.example.org,piano,2
EOT

    bin/dataset -quiet -nl=false mydata.ds import-csv freda.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds import-csv freda.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false mydata.ds attach freda freda.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds attach freda freda.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false mydata.ds import-csv mojo.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds import-csv mojo.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false mydata.ds attach mojo mojo.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds attach mojo.csv'
        exit 1
    fi
    bin/dataset -quiet -nl=false mydata.ds attachments mojo
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds attachments mojo'
        exit 1
    fi
    if [[ -f "mojo.csv" ]]; then
        rm mojo.csv
    fi
    bin/dataset -quiet -nl=false mydata.ds detach mojo mojo.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds attachments mojo mojo.csv'
        exit 1
    fi
    if [[ ! -f "mojo.csv" ]]; then
        echo 'test_attachments: (failed) mydata.ds detach mojo mojo.csv'
        exit 1
    fi
    bin/dataset -quiet -nl=false mydata.ds prune freda freda.csv
    if [[ "$?" != "0" ]]; then
        echo 'test_attachments: (failed) mydata.ds prune freda freda.csv'
        exit 1
    fi

    # Success, cleanup our test data
    rm freda.csv mojo.csv
    rm -fR mydata.ds
}

function test_check_and_repair() {
    if [[ -d "myfix.ds" ]]; then
        rm -fR myfix.ds
    fi
    cat << EOT > myfix.csv
Name,EMail,Office,Count
freda,freda@inverness.example.edu,4th Tower,1
mojo,mojo.sam@sams-splace.example.org,piano,2
EOT
    bin/dataset -quiet -nl=false myfix.ds init
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) myfix.ds init'
        exit 1
    fi
    bin/dataset -quiet -nl=false myfix.ds import-csv myfix.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) myfix.ds import-csv myfix.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false myfix.ds check
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) myfix.ds check'
        exit 1
    fi
    echo '{}' > myfix.ds/collection.json
    bin/dataset -quiet -nl=false myfix.ds check
    if [[ "$?" != "1" ]]; then
        echo 'test_check_and_repair: (failed, expected exit code 1) myfix.ds check'
        exit 1
    fi
    bin/dataset -quiet -nl=false myfix.ds repair
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) myfix.ds repair'
        exit 1
    fi
    bin/dataset -quiet -nl=false keys
    if [[ "$?" != "0" ]]; then
        echo 'test_check_and_repair: (failed) myfix.ds repair'
        exit 1
    fi
   
    # Success, cleanup
    rm -fR myfix.ds
    rm myfix.csv
}

function test_count() {
    echo 'Test dataset count'
    if [[ -d "count.ds" ]]; then
        rm -fR count.ds
    fi
    cat << EOT > count.csv
Name,EMail,Office,Count,published
freda,freda@inverness.example.edu,4th Tower,1,true
mojo,mojo.sam@sams-splace.example.org,piano,2,false
EOT
 
    if [[ ! -f "count.csv" ]]; then
        echo 'test_count: (failed) could not create count.csv'
        exit 1
    fi

    bin/dataset -quiet -nl=false count.ds init
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) count.ds init'
        exit 1
    fi
    bin/dataset -quiet -nl=false count.ds import-csv count.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) count.ds import-csv count.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false count.ds count
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) count.ds count'
        exit 1
    fi
    bin/dataset -quiet -nl=false count.ds count '(eq .published true)'
    if [[ "$?" != "0" ]]; then
        echo 'test_count: (failed) count.ds count "(eq .published true)"'
        exit 1
    fi

    # Success, cleanup
    rm -fR count.ds
    rm count.csv
}

function test_search() {
    echo 'Testing indexing and find'
    if [[ -d "search.ds" ]]; then
        rm -fR search.ds
    fi
    cat << EOT > search.csv
id,title,type,published,author
a1,4th Tower of Inverness,audio play,true,Tom Lopez
n1,"20,000 leagues under the Sea",novel,true,Jules Verne
s1,Our Person in Avalon,screenplay,false,R. S. Doiel
EOT
    
    cat << EOT > search.json
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

    bin/dataset -quiet -nl=false search.ds init 
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) search.ds init'
        exit 1
    fi
    bin/dataset -quiet -nl=false search.ds import-csv search.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) search.ds import-csv search.csv 1'
        exit 1
    fi
    bin/dataset -quiet -nl=false search.ds indexer search.json search.bleve
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) search.ds indexer search.json search.bleve'
        exit 1
    fi
    bin/dataset -quiet -nl=false find search.bleve 'screenplay' 
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) search.ds find "screenplay"'
        exit 1
    fi
    bin/dataset -quiet -nl=false find search.bleve '+published:true' 
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) search.ds find "+published:true"'
        exit 1
    fi
    bin/dataset -quiet -nl=false find search.bleve 'Tower'
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) search.ds find "Tower"'
        exit 1
    fi
    echo 'a1' > list.keys
    bin/dataset -quiet -nl=false deindexer search.bleve list.keys
    if [[ "$?" != "0" ]]; then
        echo 'test_search: (failed) deindexer search.bleve list.keys' 
        exit 1
    fi

    # Success, cleanup
    rm -fR search.ds
    rm -fR search.bleve
    rm search.json
    rm search.csv
    rm list.keys
}

function test_import_export() {
    echo 'Test import-csv, export-csv'
    if [[ -d "pubs.ds" ]]; then
        rm -fR "pubs.ds"
    fi
    cat << EOT > in.csv
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

    bin/dataset -quiet -nl=false pubs.ds init
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) pubs.ds init'
        exit 1
    fi

    bin/dataset pubs.ds "import-csv" in.csv 1
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) pubs.ds import-csv in.csv 1'
        exit 1
    fi
    bin/dataset pubs.ds keys
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) pubs.ds keys'
        exit 1
    fi
    bin/dataset pubs.ds "export-csv" "out.csv" "true" '._Key,.title,.type,.date_type,.date' 'EPrint ID,Title,Type, Date Type,Date'
    if [[ "$?" != "0" ]]; then
        echo 'test_import_export: (failed) pubs.ds export-csv out.csv true "._Key,.title,.type,.date_type,.date" "EPrint ID,Title,Type,Date Type,Date"'
        exit 1
    fi

    # Success, cleanup
    rm -fR data.ds
    rm in.csv
    rm out.csv
}

echo "Testing command line tools"
test_dataset
test_gsheet
test_issue15
test_issue19
test_readme
test_getting_started
test_attachments
test_check_and_repair
test_count
test_search
test_import_export
echo 'Success!'
