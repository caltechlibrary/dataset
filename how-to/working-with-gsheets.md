
# Working with GSheets

## A walking through the process

In this walk through we will create some data in a Google Sheet, import 
it into a collection we've created. Modify it and send it back to our 
Google Sheet. As part of the procees we will need to authorize our 
_dataset_ tool and setup access.

### Setup a local collection

Setup our local collection. Create a new dataset collection named 
"ZBS-Character-List.ds"

```shell
    dataset init zbs-cast-list.ds
```

Download [zbs-cast-list.csv](zbs-cast-list.csv).  We will use this when 
we create our Google Sheet for this walk through.

### Setting up our Google Sheet

Open your browser and create a new Google Sheet by going to 
https://sheets.google.com. You can pick the "Blank" sheet from the 
template gallery under "Start a new spreadsheet".  This should you a 
new untitled spreadsheet.  Set the title to something meaninful, I'm 
going to set my title to "zbs-cast-list".

### loading some sample data

In Google Sheets go to the file menu and select "import".  Select the 
"upload" tab, find _zbs-cast-list.csv_ and "open" it. You should then see 
an "Import file" dialog, select "Replace current sheet" then press 
"Import Data" button at the bottom of the box.

You should wind up with a spreadsheet that starts out something like this --

```csv
    ID,Name,Title,Year
    1,Jack Flanders,The Fourth Tower of Inverness,1972
    2,Little Freida,The Fourth Tower of Inverness,1972
    3,Dr. Mazoola,The Fourth Tower of Inverness,1972
    4,The Madonna Vampyra,The Fourth Tower of Inverness,1972
    5,Chief Wampum,The Fourth Tower of Inverness,1972
    6,Old Far-Seeing Art,The Fourth Tower of Inverness,1972
    7,Lord Henry Jowls,The Fourth Tower of Inverness,1972
    8,Meanie Eenie,The Fourth Tower of Inverness,1972
    9,Lady Sarah Jowls,The Fourth Tower of Inverness,1972
    ...
```

You sheet should have four columns A to D (ID, Name, Title, Year) and 
195 rows with id ranging from 1 to 194 in column 1. You can describe the 
range of your sheet as A1:D195 (A 1 through D 195). We're going to use 
this range when importing and exporting.  

### Finding your sheet's ID and sheet name

Look at the URL for your Google Sheet. It'll be similar to mine

```
    https://docs.google.com/spreadsheets/d/1Rf-fbGg9OPWnDsWng9oyQmMIWzMqx717uuKeBlzDaCc/edit#gid=0
                                           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
```

Notice the part of the URL between "/d/" and "/edit". For my sheet it is 
"1Rf-fbGg9OPWnDsWng9oyQmMIWzMqx717uuKeBlzDaCc". This is the sheet ID. It 
is unique to my spreadsheet so yours will be different. We are going to 
need that number later.

Look at the bottom of the Google Sheet window. You'll see the tab for 
your sheet. By default this is named "Sheet1", since I chose to replace my 
current empty sheet it named my sheet "zbs-cast-list" as that was the name 
of the file I imported.  You'll need to know the sheet name as well as 
they sheet ID.

### Setting up authorization to access your Google Sheet

This is the hardest part of the process. Google sometimes updates this
procedure, and this worked as of 08/08/2018. You might have to make 
adjustments in the future.

We need to get some credentials and authorization keys to access your 
Google Sheet before we can import it into our dataset collection.  If 
you have not done this previously go to 
https://developers.google.com/sheets/api/quickstart/go. We only need to 
complete "Step 1". 

Click the blue button labeled `Enable the Google Sheets API`

Make up a project name (e.g. zbs-cast-list).

Click the blue button labeled `Download Client Configuration`

Currently, you need to move the file from your downloads folder to your
current working directory (e.g. on my Mac I did 
`mv ~/Downloads/credentials.json ./`

### Authorizing dataset access

Once your "credentials.json" file is generated and downloaded to your 
working directory you need to trigger authorization for _dataset_. To 
do so we type in for importing our spreadsheet into our local collection.  
The first time we do this a link (URL) will be displayed. Copy that link 
into your web browser. You will goto a page that allows you to "authorize" 
the application.  Follow the instructions.

    NOTE: The "credentials.json" file and OAuth authorization is 
    required for data to access your Google Sheet.  The OAuth 
    authorization process uses the "credentials.json" file to get 
    a token to use on subsequent access by _dataset_. This OAuth 
    token is usually stored in your `$HOME/.credentials` directory as 
    sheets.googleapis.com-dataset.json.  If this file doesn't exist 
    then the first time you run the _dataset_ command with a GSheet 
    option it'll prompt you to use your web browser to authorize 
    _dataset_ to access your Google spreadsheet.

The command that will trigger the authorization process the first time is 
also the command we will eventually use to import our data. Replace 
SHEET_ID with the number we see in the URL and SHEET_NAME with the name 
of the spreadsheet. Our CELL_RANGE will be "A1:D195" and our 
COL_NP_FOR_ID will be 1

```shell
    dataset import zbs-cast-lis.ds SHEET_ID SHEET_NAME COL_NO_FOR_ID [CELL_RANGE]
```

For my URL, SHEET_ID, SHEET_NAME, CELL_RANGE and COL_NO_FOR_ID looked like

```shell
    dataset import zbs-cast-lis.ds \
       "1Rf-fbGg9OPWnDsWng9oyQmMIWzMqx717uuKeBlzDaCc" \
       "zbs-cast-list" \
       1 "A1:D192"
```

Yours will have a different SHEET_ID and SHEET_NAME.

After your authorize the sheet access via your web browser then next 
time you run the command you'll see data imported into _zbs-cast-list.ds_. 
You can count the keys to see what was imported.

```shell
    dataset zbs-cast-list.ds count
```

You are now ready to modify and update your local collection.

## Synchronization with a Google Sheet

_dataset_ provides a mechanism to synchronize data with a table (e.g. 
a CSV file or Google Sheet, we're interested in the later).  
Synchronizationis accomplished by mapping the column headings to object 
paths in our dataset collection as well as rows to objects. To define 
this relationship so dataset knowns what to do we use a "data frame". 
For dataset this means describing a set of keys (which will map to rows 
in the table), a set of dotpaths (this well become columns mapped into 
each row) and labesl (the names of the columns in the table form). Let's 
create a "frame" for synchronizing zbs-cast-list.ds with our Google Sheet.

First we need to make sure we know our fields in our JSON objects.
We can see a random sample using the keys verb and retrieving the
resulting key list as a list of objects and pretty printing them.

```shell
    dataset keys -sample=1 zbs-cast-list.ds | \
        dataset read -p -i - zbs-cast-list.ds
```

Here is an example of the output

```json
    {
        "_Key": "19",
        "ID": 19,
        "Name": "Comtese Zazeenia",
        "Title": "Moon Over Morocco",
        "Year": 1973
    }
```

Notice that we a `_Key` field and an "ID" field. `_Key` is the internal
id for the JSON object used by dataset. We will want to
use the `_Key` field explicitly when we defined our frame. This
will establish the relationship between dataset's objects and
the spreadsheet.  We're going to want to include "all" keys in the
collection so we'll be using the '-all' option (you could
limit the frame to specific records by providing a keylist 
to the frame definition).

Step 1. define our frame

```shell
    dataset frame-create -all zbs-cast-list.ds gsheet-sync \
        ._Key=ID .Name=Name .Title=Title .Year=Year
```

This returns a new frame definition. This includes the relationship
between our object attributes (dot paths) and the column label.

Step 2. Review the frame you defined.

```shell
    dataset frame -p zbs-cast-list.ds gsheet-sync | more
```

You should check your recreate the frame and if the dot paths or 
labels look incorrect. Otherwise we're ready to synchronize our 
collection with our Google Sheet.

Let's change item 43 in our Google Sheet from "Jack Flanders" to
"Molly Flanders". We want our collection to pick up this change. We need
to "recieve" data into our collection from our Google Sheet. We need
to do a "sync-receive".

```shell
    dataset sync-recieve zbs-cast-list.ds gsheet-sync \
        1tXbMC1Dt5B8sFr1MvJAuwkS3TatVssu0f4YcAJoZgOE \
        zbs-cast-list
```

We can check that we recieved our data by "reading" the record 43
in our collection.

```shell
    dataset read -p zbs-cast-list.ds 43
```

Molly has been promoted to director and Jack is back in the cast.
Let's update our collection then "send" our data back to the Google Sheet.

```shell
    dataset join -overwrite zbs-cast-list.ds 43 '{"name":"Jack Flanders"}'
```

Now let's update send our frame back up to our Google Sheet.

```shell
    dataset sync-send zbs-cast-list.ds gsheet-sync \
        1tXbMC1Dt5B8sFr1MvJAuwkS3TatVssu0f4YcAJoZgOE \
        zbs-cast-list
```

Technically all the rows/columns in our Google Sheet we updated.
If we changed our frame to only have key 43 then it would have only updated
the row with the matching ID of 43.

## Exporting a collection to a Google Sheet

We can also export our collection to Google Sheet. Exporting overwrites
any content in the sheet with our collections' frame's column and row 
order. "sync-send" respects the existing spreadsheet column and row
order, export imposes a the collection frame's column and row order.

This command to export our collection into Google sheets looks like

```shell
    dataset export COLLECTION_NAME FRAME_NAME SHEET_ID SHEET_NAME 
```

Like "sync-receive" and "sync-send" we use frame to define our export.
In our next example we're create a new sheet called "new-cast-list" in
our Google Sheet, then we can export our whole zbs-cast-list.ds into it.

```shell
    dataset export zbs-cast-list.ds gsheet-sync \
        1tXbMC1Dt5B8sFr1MvJAuwkS3TatVssu0f4YcAJoZgOE \
        new-cast-list
```

You should now see populated new-cast-list sheet.


Related topics: [dotpath](../docs/dotpath.html), [export-csv](../docs/export-csv.html), [frame](../docs/frame.html), [import-csv](../docs/import-csv.html), [import-gsheet](../docs/import-gsheet.html), [export-gsheet](../docs/export-gsheet.html), [sync-receive](../docs/sync-receive.html) and [sync-send](../docs/sync-send.html) 

