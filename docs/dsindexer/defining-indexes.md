# Defining Indexes

_dsindexer_ builds an index from an index map file.  A map defines the structure of the index. The definition file is a JSON document.
_dsindexer_ supports two types of map files. A simple version and also the more complicated version native to the Bleve search package.

## The Simple index map

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


FIXME: This is expeculation on how defining complex indexes might work.

## Indexing more complex JSON documents

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
In our simple approach we could describe the title and three authors explicitly like this.

```json
   {
       "title": {
           "object_path": ".title"
       },
       "author_1":{
           "object_path": ".authors[0].sort_name"
       },
       "author_2": {
           "object_path": ".authors[1].sort_name"
       },
       "author_3": {
           "object_path": ".authors[2].sort_name"
       }
   }
```

The trouble is what if we want to index display name and sort name independantly? What if we have 100 authors instread of three.
This simple approach of explicit paths quickly becomes problematic. What we need to do is beable to describe to Bleve how to reach
into our tree and pull out the pieces we're interested in. It's a problem of notation really. If your writing a custom indexer in
Go the Bleve package has functions for handling but this leaves us with the problem of how do we easily describe in our
definition file those more complex relationships?

The approach _dataset_ takes when describing the index structure is to nest the definitions just like the data structure we're
describing. Let's take another pass at describing our article metadata.


index can reach into

```json
    {
       "title": {
           "object_path": ".title"
       },
       "authors_display_name": {
            "object_path": ".authors[:].display_name"
       },
       "authors_sort_name": {
            "object_path": ".authors[:].sort_name"
       },
       "authors_orcid": {
            "object_path": ".authors[:].orcid"
       }
    }
```

Notice that we've create an array os the value for "authors".  In the array we have a single object that describes what the array
is holding. If we're working with an array objects then an anonymous object is described with each property of the object
named and defined with a dot path in relationship to the object. If we were describing an array of strings we'd still describe
it with an anonymous object but the dotpath would only contain a single period "." as its relative root.

Here's an example where what an array of years might look like as a definition

```json
       "years": {
          "object_path": ".years[:]",
          "field_mapping": "numeric"
       }
    }
```

_dsindexer_ will only index arrays that containing a single data type.  So if you have an array that has an object,
a numeric value and a string you're out of luck or you'll need to index each type separately.



## The Bleve native index map

_dsindexer_ works from a index definition expressed as a JSON document. It is the same format as Bleve's native
index definition in JSON. Bleve native indexes are distinguished by the file extension `.bmap`.  Bleve supports complex 
including things like facetted search.  In our example we'll keep it simple indexing only two specfic fields -- 
last_name ad date_of_birth.

If your JSON data document looks like

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

You could create an index of last name and date of birth (e.g. `last_name-dob.bmap`)  
with the following definition

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

