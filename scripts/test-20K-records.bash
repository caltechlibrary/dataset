#!/bin/bash

function mkRec() {
    K="${1}"
    T="${2}"
    cat <<EOF
{"_id":"${K}","t":"${T}"}
EOF
}

COLLECTION_NAME="TestCollection-$(date +%Y%m%d)"
if [ ! -d "${COLLECTION_NAME}" ]; then
    dataset init "${COLLECTION_NAME}"
fi
export DATASET="${COLLECTION_NAME}"

echo "Building 20K records in ${DATASET}"
for I in $(range 1 20000); do
    T="$(date)"
    REC="$(mkRec "${I}" "${T}")"
    bin/dataset -c "${COLLECTION_NAME}" create "$I" "${REC}"
    echo "${REC}"
    I=""
    T=""
done

C=$(bin/dataset count)
I=$(bin/dataset keys | wc -l)
if [ "$C" != "$I" ]; then
    echo "Expected $C, got $I"
    exit 1
fi
if [ -d "$COLLECTION_NAME" ]; then
    rm -fR "${COLLECTION_NAME}"
fi
echo "Successful test run"
