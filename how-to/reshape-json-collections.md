Reshape JSON collections
========================

Problem
-------

You have a __dataset__ collection of JSON documents but the fields you\'re
interested in are nested. Using __dataset__ and
[datatools](https://github.com/caltechlibrary/datatools/releases/latest)\'s
__jsonmunge__ you can reshape your existing collection\'s record to the
shape you prefer.

In this how to we will look at mondify a single record then once we have
the record looking the way we want apply that transformation to the
whole collection.

Example
-------

In our collection we have record with an id of 12345. Running
`dataset read 12345` we can see our record looks like\--

``` {.json}
    {
        "title": "The wonderful world of data collecting",
        "authors": [
            {"family": "Brown", "given": "Jules"},
            {"family": "Brown", "given": "Verne"}
        ]
    }
```

What we\'d like is a flattened version of the author names.

``` {.json}
    {
        "title": "The wonderful world of data collecting",
        "author_display_names": "Jules Brown and Verne Brown",
        ...
    }
```

We\'re going to pull out each others name object and then format them
the way we prefer. __jsonmunge__ lets us apply a Go text template to our
JSON data and then output something. In our case our formatted names.

### Building up a template

Generating our `.author_display_names` field can be broken down into
simpler parts. First we are going to look at formatting a single name
and then look at how to format both names and finally format an number
of names. Inside `.authors` array we have a name object. It has
`.family` and `.given` attributes. A simple template would reach in to
the `.authors` array by index and then order the `.given` and `.family`
attributes as desired. Array indexes count from zero so the first
author\'s index is zero. The template function *dotpath* lets us reach
inside the array.

Try this

``` {.shell}
    dataset read 12345 | \
       jsonmunge -i - -E '{{ dotpath . ".authors[0].given" "" }} {{ dotpath . ".authors[0].family" "" }}'
```

Let\'s take this command pipeline apart. We retrieved our dataset record
12345 with `dataset read 12345`. We send that record to __jsonmunge__
(`-i -` is idiomatic of datatool commands for saying read from standard
input since the record should be coming from __dataset__\'s standard
output) and the `-E` to evaluate a simple template ordering out first
author name.

        Jules Brown

It\'s a bit ugly (and long) but we can adapt that to display both names.

``` {.shell}
    dataset read 12345 | \
       jsonmunge -i - -E '{{ dotpath . ".authors[0].given" "" }} {{ dotpath . ".authors[0].family" "" }} and {{ dotpath . ".authors[1].given" "" }} {{ dotpath . "authors[1].family" "" }}'
```

getting

        Jules Brown and Verne Brown

That command line is getting pretty long. Let\'s take that expression
and put it in a template file called \"flatten.tmpl\".

        {{ dotpath . ".authors[0].given" "" }} {{ dotpath . ".authors[0].family" "" }} and {{ dotpath . ".authors[1].given" "" }} {{ dotpath . "authors[1].family" "" }}

Run the template and see the results with

``` {.shell}
    dataset read 12345 | jsonmunge -i - flatten.tmpl
```

We should again see

        Jules Brown and Verne Brown

What happens for the next record where the number of authors is
different? Looking at our original data we see that `.authors` is an
array of objects. Go\'s text templates have a function called *range*
which makes it easy to iterate over arrays. *range* can return the index
value as well as the object at that index. Applying the *range* function
would look like this version of \"flatten.tmpl\".

        {{ range $i,$author := .authors }}
            {{ if (gt $i 0) }} and {{ end }}
            {{ $author.given }} {{ $author.family }}
        {{ end }}

Running

``` {.shell}
    dataset read 12345 | jsonmunge -i - flatten.tmpl
```

we get

           Jules Brown

            and
           Verne Brown

That sorta gives us what we wanted but the spacing is all wrong and we
have some extra line breaks. We could put all the template parts in one
line but that would make it hard to read and debug. Fortunately Go
templates elements and indicate if leading or trailing whitespace should
be trimmed. You do that by using `{{-` and `-}}` for trimming leading
and trailing whitespace. The revised template will look like

        {{- range $i,$author := .authors }}
            {{- if (gt $i 0) }} and {{ end -}}
            {{- $author.given }} {{ $author.family -}}
        {{- end -}}

and a running that through
`dataset read 12345 | jsonmunge -i - flatten.tmpl` gives us

        Jules Brown and Verne Brown

Ok, so how does this help us reshape our origin 12345 record? Well first
we need to turn our string \"Jules Brown and Verne Brown\" into an
object. Updating our template the curly brackets and attribute nations
gives us

        {
            "author_display_names": "{{- range $i,$author := .authors }}
                {{- if (gt $i 0) }} and {{ end -}}
                {{- $author.given }} {{ $author.family -}}
            {{- end -}}"
        }

Now running `dataset read 12345 | jsonmunge -i - flatten.tmpl` gives us
our new object.

        {
            "author_display_names": "Jules Brown and Verne Brown"
        }

Now we ready to \"join\" our new object with the 12345 record. We can do
that by extending our pipe line.

``` {.shell}
    dataset read 12345 | jsonmunge -i - flatten.tmpl | dataset -i - join update 12345
```

We can check to make sure it worked with `dataset read 12345`. You
should see something like (order of attributes may vary)

``` {.json}
    {
      "author_display_names": "Jules Brown and Verne Brown",
      "authors": [
        { "family": "Brown", "given": "Jules" },
        { "family": "Brown", "given": "Verne" }
      ],
      "title": "The wonderful world of data collecting"
    }
```

Notice we now have a new `.author_display_names` attribute in our
object. We still see our old `.authors`. The *join* function will not
overwriting fields nor trim others. It adds the attributes of one object
to the other.

Now let\'s say you decide you\'d rather have the names in \"family,
given\" order for the individual names. Using *join -overwrite* we can
replace the value in `.author_display_names` with a new one.

**flatten.tmpl** should now look like

        {
            "author_display_names": "{{- range $i,$author := .authors }}
                {{- if (gt $i 0) }} and {{ end -}}
                {{- $author.family -}}, {{ $author.given -}}
            {{- end -}}"
        }

Running

``` {.shell}
    dataset read 12345 | \
      jsonmunge -i - flatten.tmpl | \
          dataset -i - join overwrite 12345` 
```

yields our new results

``` {.json}
    {
        "title": "The wonderful world of data collecting",
        "authors": [
            {"family": "Brown", "given": "Jules"},
            {"family": "Brown", "given": "Verne"}
        ],
        "author_display_names": "Brown, Jules and Brown, Verne"
    }
```

### Putting it together

We can transform a single record but how about transforming the entire
collection? That turns out to be easy we just loop over each key in the
collection applying our pipeline.

``` {.shell}
    dataset keys | while read K; do
        dataset read "$K" | \
           jsonmunge -i - flatten.tmpl | \
               dataset join overwrite "$K"
    done
```

Where we had \"12345\" before we now have `"$K"`. The rest is just
waiting on the computer to finish.
