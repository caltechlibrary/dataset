VERB

sync-send

    sync-send COLLECTION FRAME_NAME [CSV_FILENAME|GSHEET_ID SHEET_NAME [CELL_RANGE]]

sync a frame of objects sending data to a table (e.g. CSV, GSheet)

OPTIONS

    -O, -overwrite  overwrite existing cells in table
    -client-secret  (sync-send to a GSheet) set the client secret path and filename for GSheet access
    -i, -input  read CSV content from a file
    -o, -output  write CSV content to a file
    -v, -verbose  verbose output


