#!/bin/bash

# Source the configuraiton file with the
# NOAA API access Token, see: https://www.ncdc.noaa.gov/cdo-web/token
# To get an access token
if [[ -f "etc/noaa-api-v.bash" ]]; then
    . etc/noaa-api-v2.bash
fi

# Make sure the NOAA_ACCESS_TOKEN is set properly
if [[ "${NOAA_ACCESS_TOKEN}" == "" ]]; then
    echo "This demo requires NOAA_ACCESS_TOKEN, see https://www.ncdc.noaa.gov/cdo-web/token"
    exit 1
fi

# Make sure we have a dataset collection, if not create the collection
dataset status noaa-demoset.ds
if [[ "$?" != "0" ]]; then
    dataset init noaa-demoset.ds
fi
export DATASET="noaa-demoset.ds"

NOAA_API_URL="https://www.ncdc.noaa.gov/cdo-web/api/v2"

# Get NOAA Location data for Saugus, CA, US.
LOCATION_NAME="SAUGUS CALIFORNIA, CA US"
STATION_ID="GHCND:USR0000CSAU"
echo "Harvesting the station ($STATION_ID) for ${LOCATION_NAME}"
ENDPOINT="stations"
curl -s -H "token:${NOAA_ACCESS_TOKEN}" "${NOAA_API_URL}/${ENDPOINT}/${STATION_ID}" | dataset -i - create "saugus-${ENDPOINT}"
# Get NOAA datasets for Saugus station
ENDPOINT="datasets"
echo "Getting ${ENDPOINT} for ${STATION_ID}"
curl -s -H "token:${NOAA_ACCESS_TOKEN}" "${NOAA_API_URL}/${ENDPOINT}?statusid=${STATION_ID}" | dataset -i - create "saugus-${ENDPOINT}"
# Get NOAA Datasets for Saugus station
#       We will loop throught the datasets available for Saugus station
YEAR="2001"
ENDPOINT="data"
echo "Getting ${ENDPOINT} ($YEAR) for ${STATION_ID}"
dataset read saugus-datasets | jsoncols -i - '.results[:].id' | jsonrange -i - -values | sed -E 's/"//g' | while read DATASETID; do
    echo "Getting ${ENDPOINT} for ${STATION_ID} dataset $DATASETID in ${YEAR}"
    PARAMS="datasetid=${DATASETID}&stationid=${STATION_ID}&startdate=${YEAR}-01-01&enddate=${YEAR}-12-31&limit=1000"
    curl -s -H "token:${NOAA_ACCESS_TOKEN}" "${NOAA_API_URL}/${ENDPOINT}/${DATACATEGORY}?${PARAMS}" | dataset -i - create "saugus-${DATASETID}-${YEAR}"
done
echo -n "Total records harvested"
dataset count
echo "Available records harvested:"
dataset keys
echo "Saugus station"
#{"elevation":442,"mindate":"1994-09-01","maxdate":"2017-12-01","latitude":34.425,"name":"SAUGUS CALIFORNIA, CA US","datacoverage":0.9962,"id":"GHCND:USR0000CSAU","elevationUnit":"METERS","longitude":-118.525}
SRC=$(dataset read saugus-stations)
cat <<EOF

    Name: $(echo "${SRC}" | jsoncols -i - '.name')
    Station ID: $(echo "${SRC}" | jsoncols -i - '.id')
    Lat/Long: $(echo "${SRC}" | jsoncols -i - '.latitude' '.longitude')
    Elevation: $(echo "${SRC}" | jsoncols -i - '.elevation' '.elevationUnit' | sed -E 's/,/ /g')

EOF

