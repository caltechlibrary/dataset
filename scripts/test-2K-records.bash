#!/bin/bash

TEST_SIZE=2000

if [ "$1" != "" ]; then
    TEST_SIZE="$1"
fi

function mkRec() {
	K="${1}"
	T="${2}"
	cat <<EOF
{"_id":"${K}","t":"${T}"}
EOF
}

echo "Setting up test collections"
COLLECTION_1="TestCollection-01-$(date +%Y%m%d)"
COLLECTION_2="TestCollection-02-$(date +%Y%m%d)"
if [ ! -d "${COLLECTION_1}" ]; then
	echo "Initializing $COLLECTION_1"
	bin/dataset init "${COLLECTION_1}"
fi
if [ ! -d "${COLLECTION_2}" ]; then
	echo "Initializing $COLLECTION_2"
	bin/dataset init "${COLLECTION_2}"
fi
export DATASET="${COLLECTION_1}"

echo "Creating 2K records in ${DATASET}"
for I in $(range 1 "$TEST_SIZE"); do
	T="$(date)"
	REC="$(mkRec "${I}" "${T}")"
	echo -n "Create $I -> ${REC} "
	bin/dataset -c "${COLLECTION_1}" create "$I" "${REC}"
	if [ "$?" != "0" ]; then
		echo "bin/dataset exited with non-zero status"
		exit "$?"
	fi
	I=""
	T=""
done

C1=$(bin/dataset count)
I1=$(bin/dataset keys | wc -l | sed -E 's/ //g')
if [ "$C1" != "$I1" ]; then
	echo "Counting keys in $DATASET, Expected $C1, got $I1"
	exit 1
fi
unset DATASET

#
# Test Copying from collection 1 to 2
#
echo "Testing copying 2K records from ${COLLECTION_1} to ${COLLECTION_2}"
bin/dataset -c "${COLLECTION_1}" keys | while read K; do
	echo -n "Copying $K "
	bin/dataset -c "${COLLECTION_1}" "read" "$K" | bin/dataset -i - -c "${COLLECTION_2}" "create" "$K"
	if [ "$?" != "0" ]; then
		echo "bin/dataset exited with non-zero status"
		exit "$?"
	fi
done

echo "Checking counts for copied data"
export DATASET="${COLLECTION_2}"
C2=$(bin/dataset count)
I2=$(bin/dataset keys | wc -l | sed -E 's/ //g')
if [ "$C2" != "$I2" ]; then
	echo "Counting keys in $DATASET, expected $C2, got $I2"
	exit 1
fi
if [ "$C1" != "$C2" ]; then
	echo "Comparing collection counts, expected $C1, got $C2"
	exit 1
fi
unset DATASET

#
# Comparing records
#
echo "Comparing records between $COLLECTION_1 and $COLLECTION_2"
bin/dataset -c "${COLLECTION_1}" keys | while read K; do
    echo -n "Comparing records for key $K"
	REC_1=$(bin/dataset -c "${COLLECTION_1}" "read" "$K")
	REC_2=$(bin/dataset -c "${COLLECTION_2}" "read" "$K")
	if [ "${REC_1}" != "${REC_2}" ]; then
		echo "Comparing records $K, expected ${REC_1} got ${REC_2}"
		exit 1
	fi
    echo " OK"
done

#
# Test list functionality
#
echo -n "Testing 'dataset list' for keys coming from stdin "
C1=$(range 1 "$TEST_SIZE" | tr " " "\n" | bin/dataset -i - -c "${COLLECTION_2}" list | jq '. | length')
if [ "$C1" != "$TEST_SIZE" ]; then
    echo "List should return 10 items, got $C1"
    exit 1
fi
echo " OK"


if [ -d "$COLLECTION_1" ]; then
	rm -fR "${COLLECTION_1}"
fi
if [ -d "$COLLECTION_2" ]; then
	rm -fR "${COLLECTION_2}"
fi
echo "Successful test run"
