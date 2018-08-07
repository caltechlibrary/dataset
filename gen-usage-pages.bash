#!/bin/bash

if [[ ! -f bin/dataset ]]; then
	echo "Running make before generating updated docs."
	make
fi
bin/dataset -generate-markdown >"docs/dataset.md"

for TOPIC in "attach" "attachments" "check" "clone" "clone-sample" "count" "create" "deindexer" "delete" "delete-frame" "detach" "export-csv" "export-gsheet" "find" "frame" "frame-labels" "frame-types" "frames" "grid" "haskey" "import-csv" "import-gsheet" "indexer" "init" "join" "keys" "list" "path" "prune" "read" "reframe" "repair" "status" "update"; do
	bin/dataset -help "${TOPIC}" >"docs/${TOPIC}.md"
done
