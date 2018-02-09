
# Defining Indexes

_dsindexer_ builds an index from an index definition file. The definition file is a JSON document. It is the same
JSON structure as used by the [bleve](https://blevesearch.com) command line tool.

## A minimal index definition

_dsindexer_ works from a index definition expressed as a JSON document. It is the same format as Bleve's native
index definition in JSON.  Bleve supports complex indexing but in this initial one we will index two fields
specifically - last_name ad date_of_birth.

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
        "types": {
            "default": {
                "enabled": true,
                "dynamic": true,
                "fields": [
                    {
                        "name": "last_name",
                        "type": "text",
                        "analyzer": "standard",
                        "store": true,
                        "index": true
                    },
                    {
                        "name": "date_of_birth",
                        "type": "datetime",
                        "store": true,
                        "index": true
                    }
                ]
            }
        }
    }
```


### Working with field mappings

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

_dsindexer_ supports four types and analyzers. 

Types are

+ text - this is the default and is what you would use for titles
+ datetime - use this for indexing numeric values
+ boolean - use this for indexing true/value values
+ geo - use this for indexing Geo Point data

There are five general types of non-language specific analyzers

+ custom - to define a custom analyzer
+ simple - a simple text analyzer
+ standard - the standard full text analyzer (this is usually what you start with)
+ keyword - keyword analysis
+ web - web content analyzer (e.g. you might use if you had HTML embedded in a JSON property)

Bleve indexes also support languages specific analyzers. Here's below is an example of our initial
index definition with all the defaults showning.

```json
    {
        "types": {
            "default": {
                "enabled": true,
                "dynamic": true,
                "fields": [
                    {
                        "name": "last_name",
                        "type": "text",
                        "analyzer": "standard",
                        "include_in_all": true,
                        "include_term_vectors": true,
                        "include_locations": true,
                        "index": true,
                        "store": true
                    },
                    {
                        "name": "date_of_birth",
                        "type": "datetime",
                        "include_in_all": true,
                        "include_term_vectors": true,
                        "include_locations": true,
                        "index": true,
                        "store": true
                    }
                ],
                "default_analyzer": ""
            }
        },
        "default_mapping": {
            "enabled": true,
            "dynamic": true,
            "default_analyzer": ""
        },
        "type_field": "_type",
        "default_type": "_default",
        "default_analyzer": "standard",
        "default_datetime_parser": "dateTimeOptional",
        "default_field": "_all",
        "store_dynamic": true,
        "index_dynamic": true,
        "docvalues_dynamic": true,
        "analysis": {}
    }
```

### Working with analyzers

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

Example of language analyzers supported are - 

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
        "types": {
            "default": {
                "fields": [
                    {
                        "name": "author",
                        "type": "text",
                        "analyzer": "simple"
                    },
                    {
                        "name": "title",
                        "type": "text",
                        "analyzers": "standard"
                    },
                    {
                        "name": "abstract",
                        "type": "text",
                        "analyzers": "standard"
                    }
               ]
            }
        }
    }
```

If your content was in Spanish you could use the Spanish language analyzer.

```json
    {
        "types": {
            "default": {
                "fields": [
                    {
                        "name": "author",
                        "type": "text",
                        "analyzer": "simple"
                    },
                    {
                        "name": "title",
                        "type": "text",
                        "analyzers": "standard",
                        "lang":"es"
                    },
                    {
                        "name": "abstract",
                        "type": "text",
                        "analyzers": "standard",
                        "lang":"es"
                    }
                ]
            }
        }
    }
```

If knew our documents were in German we could try something like this definition--


```json
    {
        "types": {
            "default": {
                "fields": [
                    {
                        "name": "author",
                        "type": "text",
                        "analyzer": "simple"
                    },
                    {
                        "name": "title",
                        "type": "text",
                        "analyzers": "standard",
                        "lang":"es"
                    },
                    {
                        "name": "abstract",
                        "type": "text",
                        "analyzers": "standard",
                        "lang":"es"
                    }
                ]
            }
        }
    }
```

Note you can use different analyzers on different fields. 

### Additoinal configuration

This additional configuration is useful for managing the size your of your index(es) on disc as
well as impact the ammount of time it takes to index your data.

#### Storing the field values in the index

As we define the numbers of fields in our index the size the index will also grow.  If you don't need to
see the field in the results you can choose not to store it in the index.  This is done with the "store"
attribute in the field's definition. The value can be true/false.

### Include Term Vectors

You can choose to include term vectors in your index. This is set by the field property called "include_term_vectors"
and like "store" it can be either true/false.

### Include In all

"include_in_all", indicates to include any composite fields named "_all", defaults to true, if you don't need this and
would like to make the index slightly smaller then you could set this to false.


### Date Format

The "date_format" string is used to indentify how to parse the date. The formatting pattern is based on Go's time.Parse()
module. You can read more about that here at https://golang.org/pkg/time/#pkg-constants. If you're using the "datetime"
field mapping for a field you should probably set the "date_format" too since dates can be written so many ways.


## Indexing more complex JSON documents

FIXME: This needs to be updated to show how to define and index sub-documents

One of the reason JSON is used for serialization of data is that it can represent many of the common types
of data structures in addition to primitive data types like string and number.  We've already seen how to
work with simple JSON structures as an object. The JSON object (or map) presents data as a series
of key and value pairs.  Another common data structure represented in JSON is that of an array. An
array can be thought of as a list containing some other data types. An array often contains strings or
numbers but it can also contain objects and other arrays.  In this way JSON documents can describe the
relatationship between say an article, it's title and the authors who wrote it. It can even describe
unique identifiers for authors as well as variation of their names. Here's an example

```json
    {
        "title": "Analysis of literary dog commentary of Summer '17",
        "abstract": "Bark, yip, gur, wine, Bark. That's why you said yesterday.",
        "authors": [{
            "display_name": "R. S. Doiel",
            "species": "human",
            "sort_name": "Doiel, Robert",
            "orcid": "0000-0003-0900-6903"
        },
        {
            "display_name":"Wesneday",
            "sort_name":"A Dog, Wedneday",
            "species":"canine"
        },
        {
            "display_name":"Dodger",
            "sort_name":"Daschund, Dodger",
            "species":"canine"
        }],
        "years":[
            1992,
            1998,
            2002
        ]
    }
```

I this data example we have three authors along two fields about an article written by two canines and a human.
In our simple approach we could describe the title and authors like this.

```json
   {
       "types": {
           "default": {
               "fields": [
                    {
                        "name": "sort_name",
                        "type": "text",
                        "analyzer": "standard"
                    },
                    {
                        "name": "display_name",
                        "type": "text",
                        "analyzer": "standard"
                    }
               ]
           }
       }
   }
```

The trouble is what if we want to index behavior display name and sort name to be independant (e.g. treat sort_name more like a keyword)? 
We can do that by choosing an different analyzer from the standard one for sort_name. Bleve supports several types of anlayzers
(e.g. simple, standard, keyword, web and custom).

```json
   {
       "types": {
           "default": {
               "fields": [
                    {
                        "name": "sort_name",
                        "type": "text",
                        "analyzer": "keyword"
                    },
                    {
                        "name": "display_name",
                        "type": "text",
                        "analyzer": "standard"
                    }
               ]
           }
       }
   }
```

What about dates? In our record we have an array of years.  We can use a different "type" when defining how we want to index years.

```json
   {
       "types": {
           "default": {
               "fields": [
                    {
                        "name": "years",
                        "type": "datetime"
                    }
               ]
           }
       }
   }
```

Indexes themselves can be defined fairly simple as we have so far and aggregated together after the fact. In addition to data shapping
approaches dsindexer supports the full Bleve index functionality, see [Bleve](https://blevesearch.com).

