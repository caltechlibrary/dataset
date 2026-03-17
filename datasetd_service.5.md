%datasetd(5) user manual | version 2.3.4-test 78e4967
% R. S. Doiel and Tom Morrell
% 2026-03-16


# datasetd Service

The datasetd based application can be configured to be managed by
systemd. You need to create a an appropriate service file with
Unit, Service and Install described.

## Example

Below is a generic datasetd systemd style service file for a project
called citesearch implemented as citesearch.yaml using datasetd to provide
the web service.

~~~
[Unit]

Description=A Citation search engine for CaltechTHESIS search

[Service]
Type=simple
ExecStart=/usr/local/bin/datasetd /Sites/citesearch/citesearch.yaml

[Install]
WantedBy=multi-user.target
~~~


