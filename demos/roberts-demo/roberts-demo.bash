#!/bin/bash

unset DATASET
if [ "$1" != "" ]; then
    S3="${1}"
else
    echo "$(basename $0) S3_BUCKET_URL"
    echo "E.g. s3://test-dataset.library.example.edu/noaaweather"
    exit 
fi

# Initialize the collection and set the DATASET environment variable
echo "Initializing collection ${S3}"
E=$(dataset init ${S3})
if [ "$E" = "" ] || [ "${E:0:6}" != "export" ]; then
    echo "Problem initializing ${S3}"
    if [ "$E" != "" ]; then
        echo "Error: ${E}"
    fi
    exit 1
fi
echo "Trying '${E}'"
$E
#export DATASET="${S3}"
if [ "${DATASET}" != "${S3}" ]; then
	echo "Something went wrong DATASET not set."
	exit 1
fi
echo "Using ${DATASET}"

# Fetch some weather info
echo "Getting some weather info from NOAA..."
curl -L -o pasadena-ca-weather-codes.html 'http://www.nws.noaa.gov/nwr/coverage/ccov.php?State=CA'
curl -L -o pasadena-ca-forecast.json 'http://forecast.weather.gov/MapClick.php?lat=34.1478&lon=-118.1445&unit=0&lg=english&FcstType=json'
curl -L -o pasadena-ca-forecast.xml 'http://forecast.weather.gov/MapClick.php?lat=34.1478&lon=-118.1445&unit=0&lg=english&FcstType=dwml'

if [ ! -f pasadena-ca-forecast.json ]; then
    echo "JSON file failed to download!"
    exit 1
fi

echo "Saving pasadena-ca-forecast.json as pasadena-ca to ${S3}"
dataset create "${DATASET}" pasadena-ca pasadena-ca-forecast.json 
echo "Attaching other data files: pasadena-ca-weather-codes.html pasadena-ca-forecast.xml"
dataset attach "${DATASET}" pasadena-ca pasadena-ca-weather-codes.html pasadena-ca-forecast.xml

echo "Removing downloaded files"
/bin/rm pasadena-ca-weather-codes.html pasadena-ca-forecast.json pasadena-ca-forecast.xml

echo "Reading back new record"
dataset read "${DATASET}" pasadena-ca

echo "Listing attachments for pasadena-ca"
dataset attachments "${DATASET}" pasadena-ca

cat<<EOF

Try the following commands and see what happens in your shell

    aws s3 ls --recursive ${S3}
    dataset attachments ${S3} pasadena-ca
    dataset attached ${S3} pasadena-ca pasadena-ca-forecast.xml

EOF
