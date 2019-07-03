
# sync-send

## Syntax 

```
    dataset sync-send COLLECTION FRAME_NAME [CSV_FILENAME|GSHEET_ID SHEET_NAME [CELL_RANGE]]
```

## Description

sync a frame of objects sending data to a table (e.g. CSV, GSheet). In the case
of GSheets there is a limitation of cell size, the GSheet platform will truncate
cells if they are too long. The limitation seems to be about 50k characters.

## OPTIONS

    -O, -overwrite  overwrite existing cells in table
    -client-secret  (sync-send to a GSheet) set the client secret path and filename for GSheet access
    -i, -input  read CSV content from a file
    -o, -output  write CSV content to a file
    -v, -verbose  verbose output

Related topics: [sync-receive](sync-receive.html) [frame](frame.html)

