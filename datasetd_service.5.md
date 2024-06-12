[Unit]
Description=A Citation search engine for CaltechTHESIS search

[Service]
Type=simple
ExecStart=/usr/local/bin/datasetd /Sites/citesearch/citesearch.yaml

[Install]
WantedBy=multi-user.target

