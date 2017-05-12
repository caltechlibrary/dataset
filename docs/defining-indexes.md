
# Defining Indexes

_dsindexer_ builds an index from an index definition file. The definition file is a JSON document.

## A minimal index definition

_dsindexer_ works from a index definition expressed as a JSON document. The most important of the definition is to map
a indexed field name to a path in the JSON document being index. This is done with dotpath notation as the value associated
with a field name in the index.

If your JSON document looks like

```json
    {
        "last_name": "Smiley",
        "first_name": "George",
        "bio": {
            "date_of_birth": "1906-09-21",
            "birth_place": {
                "name": "UK, England, Salisbury",
                "geo_coord":"51.0797058,-1.8434485"
            }
        },
        "email":"smiley.george@emeritus.circus.example.net"
    }
```

You could create an index of last name and date of birth with the following definition

```json
    {
        "last_name": {
            "object_path": ".last_name"
        },
        "dob": {
            "object_path": ".bio.date_of_birth"
        }
    }
```

The dotpath notation lets you reach into a nested JSON property and bring it out into a field that will
be indexed. 

## Working with field mappings

In our example of above we have three types of data in our JSON document.  The name properties are
strings. The date of birth property is a date in YYYY-MM-DD format and finally we have an email
address. In our initial index definition we treat all these values as strings.  This is fine for
the names and email address but if we want to work with date ranges then the date of birth should
be handled differently. It should be handled as a date.

Here's a revised definition

```json
    {
        "last_name": {
            "object_path": ".last_name"
        },
        "dob": {
            "object_path": ".bio.date_of_birth",
            "field_mapping": "datetime"
        }
    }
```

_dsindexer_ supports four types of field mappings

+ text - this is the default and is what you would use for titles
+ numeric - use this for indexing numeric values
+ datetime - use this for indexing dates and time values
+ boolean - use this for indexing true/value values
+ geopoint - use this for indexing Geo Point data

If we want to expand our definition to include the location of Smiley's birth we add the geocordinates too.


```json
    {
        "last_name": {
            "object_path": ".last_name"
        },
        "dob": {
            "object_path": ".bio.date_of_birth",
            "field_mapping": "datetime"
        },
        "origin": {
            "object_path": ".bio.birth_place.geo_coord",
            "field_mapping": "geopoint"
        }
    }
```

## Working with analyzers

In addition to setting the controlling how the values are mapped into the index you can control the analysis
that are applied when building your index (see http://www.blevesearch.com/docs/Analyzers/ for details).
Analyzers include applying language rules for understanding the text analyzed. This includes handling things
like stop word removal, language settings.

_dsindexer_ support the following types of analyzers

+ keyword - performs zero analysis, use this if you want to treat the value as is
+ simple - performs minimal analysis, tokenizes using Unicode and lowercases the value
+ standard - is like simple but adds English stop word removal
+ web - tries to determine the language then applies that languages analyzer applying its rules (e.g. if
  the language detected was German then German stop words, analysis would be performed)
+ lang - will look use a language specific analyzer (relying on the lang property for language name, e.g. en, es, de, fr)

Language analyzers current supported are - 

+ Arabic (ar) 
+ Catalan (ca)
+ Chokwe (cjk) 
+ Central Kurdish (ckb) 
+ German (de)
+ English (en)
+ Spanish (es)
+ Persian (fa)
+ French (fr)
+ Hindi (hi)
+ Italian (it)
+ Portuguese (pt)

Let's consider a JSON document that has a title and abstract field.

```json
    {
        "author": "Doe, Jane",
        "title": "Some title here",
        "abstract": "blah, blah, blah, herrumph, blip, bleep"
    }
```



The default language analyzer is English (en) but you can explicitly indicate that with this definition

```json
    {
        "author": {
            "object_path": ".author",
            "field_mapping": "text",
            "analyzer": "simple"
        },
        "title": {
            "object_path": ".title",
            "field_mapping": "text",
            "analyzers": "standard"
        },
        "abstract": {
            "object_path": ".abstract",
            "field_mapping": "text",
            "analyzers": "standard"
        }
    }
```

If your content was in Spanish you could use the Spanish language analyzer.

```json
    {
        "author": {
            "object_path": ".author",
            "field_mapping": "text",
            "analyzer": "simple"
        },
        "title": {
            "object_path": ".title",
            "field_mapping": "text",
            "analyzers": "lang",
            "lang":"es"
        },
        "abstract": {
            "object_path": ".abstract",
            "field_mapping": "text",
            "analyzers": "lang",
            "lang":"es"
        }
    }
```

If knew our documents were in German we could try something like this definition--


```json
    {
        "author": {
            "object_path": ".author",
            "field_mapping": "text",
            "analyzer": "simple"
        },
        "title": {
            "object_path": ".title",
            "field_mapping": "text",
            "analyzers": "lang",
            "lang": "de"
        },
        "abstract": {
            "object_path": ".abstract",
            "field_mapping": "text",
            "analyzers": "lang",
            "lang": "de"
        }
    }
```

Note you can use different analyzers on different fields. 

## Additoinal configuration

This additional configuration is useful for managing the size your of your index(es) on disc as
well as impact the ammount of time it takes to index your data.

### Storing the field values in the index

As we define the numbers of fields in our index the size the index will also grow.  If you don't need to
see the field in the results you can choose not to store it in the index.  This is done with the "store"
attribute in the field's definition. The value can be true/false.

## Include Term Vectors

You can choose to include term vectors in your index. This is set by the field property called "include_term_vectors"
and like "store" it can be either true/false.

## Include In all

"include_in_all", indicates to include any composite fields named "_all", defaults to true, if you don't need this and
would like to make the index slightly smaller then you could set this to false.


## Date Format

The "date_format" string is used to indentify how to parse the date. The formatting pattern is based on Go's time.Parse()
module. You can read more about that here at https://golang.org/pkg/time/#pkg-constants. If you're using the "datetime"
field mapping for a field you should probably set the "date_format" too since dates can be written so many ways.


