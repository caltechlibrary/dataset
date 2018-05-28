
# Google Spreadsheet Integration

_dataset_ provides for importing from and export to a Google spreadsheet (i.e. GSheet). 
This does require some setup for this to work.  _dataset_ needs to beable to access
the Google Sheets API for reading and writing. You can find documentation on setting
up access in "step 1" at https://developers.google.com/sheets/api/quickstart/go.

You'll need a "client_secret.json" file and OAuth authorization for access to be permitted.
If credentials for the OAuth part are usually stored in your `$HOME/.credentials` directory
as sheets.googleapis.com-dataset.json.  If this file doesn't exist then the first time you
run the _dataset_ command with a GSheet option it'll prompt you to use your web browser
to authorize _dataset_ to access your Google spreadsheet. 

## FIXME: need to walk through a couple use cases so it makes sense to most people
