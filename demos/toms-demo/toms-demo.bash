#!/bin/bash

# Initialize the collection and set the DATASET environment variable
echo "Initializing collection testdataset/noaaweather.ds"
E=$(dataset init testdataset/noaaweather.ds)
$E
if [ "$DATASET" = "" ]; then
	echo "Something went wrong DATASET not set."
	exit 1
fi
echo "Using $DATASET"

# Fetch some weather info
echo "Getting some weather info from NOAA..."
curl -L -o pasadena-ca-weather-codes.html 'http://www.nws.noaa.gov/nwr/coverage/ccov.php?State=CA'
curl -L -o pasadena-ca-forecast.json 'http://forecast.weather.gov/MapClick.php?lat=34.1478&lon=-118.1445&unit=0&lg=english&FcstType=json'
curl -L -o pasadena-ca-forecast.xml 'http://forecast.weather.gov/MapClick.php?lat=34.1478&lon=-118.1445&unit=0&lg=english&FcstType=dwml'

echo "Saving pasadena-ca-forecast.json as pasadena-ca to dataset/noaaweather)"
dataset -i pasadena-ca-forecast.json create pasadena-ca
echo "Attaching other data files: pasadena-ca-weather-codes.html pasadena-ca-forecast.xml"
dataset attach pasadena-ca pasadena-ca-weather-codes.html pasadena-ca-forecast.xml

echo "Removing downloaded files"
/bin/rm pasadena-ca-weather-codes.html pasadena-ca-forecast.json pasadena-ca-forecast.xml

echo "Reading back new record"
dataset read pasadena-ca

echo "Listing attachments for pasadena-ca"
dataset attachments pasadena-ca

cat<<EOF

Try the following commands and see what happens in your shell

    ls -l
    \$(dataset init testdataset/noaaweather)
    dataset attached pasadena-ca pasadena-ca-forecast.xml
    ls -l

EOF
