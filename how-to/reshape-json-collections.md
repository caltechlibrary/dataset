
# Reshape JSON collections

## Problem

You have a dataset collection of JSON documents but the fields you're interested in are highly nested.  You can
iterate over your JSON collections and reshape with using _jsonmunge_ (from Caltech Library's 
[datatools](https://github.com/caltechlibrary/datatools/releases/latest)) and the _dataset_ *join* operation.

### Example

In our collection we have record with an id of 12345. It looks like

```json
    {
        "title": "The wonderful world of data collecting",
        "authors": [
            {"family": "Brown", "given": "Jules"},
            {"family": "Brown", "given": "Verne"}
        ]
    }
```

What we'd like is something flatter like

```json
    {
        "title": "The wonderful world of data collecting",
        "author_display_names": "Jules Brown and Verne Brown",
        ...
    }
```

There are several approaches you can take. You could copy the records from one collection to another transforming them
on the way.  You can also add the extrapolated data with back to the original collection using the *join* operation.
Its this second approach we're going to use in conjunction with _jsonmunge_.

#### Building up a template

Generating our `.author_display_names` field.  Our output should look like "given family" names for each author with an "and"
between them. But that is pretty complicated so let's just get our first name and put it in the right order (given family).

```shell
    dataset read 12345 | jsonmunge -i - -E '{{ .authors[0].given }} {{ .authors[0].family }}'
```

Let's take this command pipeline apart.  We retrieved our dataset record 12345 with `dataset read 12345`. We sent
that record to _jsonmunge_ (`-i -` is idiomatic of datatool commands for saying read from standard input) and a `-E` to
evaluate a simple template ordering out first author name.

```
    Jules Brown
```

It's a bit ugly and we could adapt that two both names like 

```shell
    dataset read 12345 | jsonmunge -i - -E '{{ .authors[0].given }} {{ .authors[0].family }} and {{ .authors[1].given }} {{ .authors[1].faimly }}'
```

getting

```
    Jules Brown and Verne Brown
```

Now that is really ugly. What happens when we have a different number of authors? Well fortunately Go's templates let you
iterate over an array with the *range* function.  *range* will return the index it is on as well as the value. Let's create
a file called "flatten.tmpl" where we can build up our new data strings. Working with the `.authors` field we'll range
over each other, put the names in order. If we're add any other beyond the first one we'll inject can "and" as needed.

```
    {{- range $i,$author := .authors -}}
        {{- if (gt $i 0) }} and {{ end -}}
        {{- $author.given }} {{ $author.family }}
    {{- end -}}
```

We've formatted our template over multiple lines and used Go template's "{{-" and "-}}" to control leading and trailing
space trimming.  If we have one author they get listed in "given family" order, if we have more (i.e. $i is more than zero)
we inject our " and ".

Running the pipeline using "flatten.tmpl" would look like

```shell
    dataset read 12345 | jsonmunge -i - flatten.tmpl
```

Our output should look like 

```
    Jules Brown and Verne Brown
```

But how do we get that back into our original JSON object as `.author_display_names`? Well as can sprinkle in some more
templating!

```
    {
        "author_display_names": "{{- range $i,$author := .authors -}}
            {{- if (gt $i 0) }} and {{ end -}}
            {{- $author.given }} {{ $author.family }}
        {{- end -}}"
    }
```

Running the pipeline again gives us the start of the JSON object we want to merge with the original record.


```shell
    dataset read 12345 | jsonmunge -i - flatten.tmpl
```

Our output should look like 

```
    {
        "author_dipslay_names":"Jules Brown and Verne Brown"
    }
```

We can update our original JSON record 12345 by sending the resulting object back to _dataset_ using the *join update* operation.

```shell
    dataset read 12345 | jsonmunge -i - flatten.tmpl | dataset -i - join update 12345
```

Now reading back the updated record

```shell
    dataset read 12345
```

with output like

```json
    {
        "title": "The wonderful world of data collecting",
        "authors": [
            {"family": "Brown", "given": "Jules"},
            {"family": "Brown", "given": "Verne"}
        ],
        "author_display_names": "Jules Brown and Verne Brown"
    }
```

The *join update* operation adds fields from the incoming object the the record in the collection. 

Now let's say you decide you'd rather have the names in "family, given" order.
Use the *join overwrite* operation after you update your template to the new form.

**flatten.tmpl** should now look like

```
    {
        "author_display_names": "{{- range $i,$author := .authors -}}
            {{- if (gt $i 0) }} and {{ end -}}
            {{- $author.family }}, {{ $author.given }}
        {{- end -}}"
    }
```

Run

```
    dataset read 12345 | jsonmunge -i - flatten.tmpl | dataset -i - join overwrite 12345
```

Now our record should looks like

```json
    {
        "title": "The wonderful world of data collecting",
        "authors": [
            {"family": "Brown", "given": "Jules"},
            {"family": "Brown", "given": "Verne"}
        ],
        "author_display_names": "Brown, Jules and Brown, Verne"
    }
```


