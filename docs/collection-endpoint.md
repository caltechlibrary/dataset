
Collection (end point)
=======================

Interacting with the __datasetd__ web service can be done with any web client. For documentation purposes I am assuming you are using [curl](https://curl.se/). This command line program is available on most POSIX systems including Linux, macOS and Windows.

The collection end point provides codemeta JSON document for a whole collection. This may include attributes for authorship, funding and contributions.

If this end point is request with a GET method then the data is returned, if requested with a POST method the date is updated the updated and a status object is returned. The POST must submit JSON encoded object conforming to the [codemeta](https://codemeta.github.io/) standards.

See the [codemeta](https://codemeta.github.io/) website for details on
creating a codemeta document.

Example
=======

The assumption is that we have __datasetd__ running on port "8485" of "localhost" and a collection named characters is defined in the "settings.json" used at launch.

Retrieving metatadata

```shell
    curl -X GET https://localhost:8485/collection/characters
```

This would return the metadata found for our characters' collection.

```json
    {
        "dataset_version": "v0.1.10",
        "name": "characters.ds",
        "created": "2021-07-28T11:32:36-07:00",
        "version": "v0.0.0",
        "author": [
            {
                "@type": "Person",
                "@id": "https://orcid.org/0000-0000-0000-0000",
                "givenName": "Jane",
                "familyName": "Doe",
                "affiliation": [
                    {
                        "@type": "Organization",
                        "@id": "https://ror.org/05dxps055",
                        "name": "California Institute of Technology"
                    }
                ]
            }
        ],
        "contributor": [
            {
                "@type": "Person",
                "givenName": "Martha",
                "familyName": "Doe",
                "affiliation": [
                    {
                        "@type": "Organization",
                        "@id": "https://ror.org/05dxps055",
                        "name": "California Institute of Technology"
                    }
                ]
            }
        ],
        "funder": [
            {
                "@type": "Organization",
                "name": "Caltech Library"
            }
        ],
        "annotation": {
            "award": "00000000000000001-2021"
        }
    }
```

