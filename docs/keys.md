
# keys

## Syntax

```
    dataset COLLECTION_NAME keys
```

## Description

List the JSON_DOCUMENT_ID available in a collection. Normally 
order is not guaranted to be the between calls. _keys_ also 
supports *filter* and *sort* expressions. For each JSON document 
which the filter expression evaluates to true for a key will be 
return.  If no sort expression is supplied the order is not 
guaranteed.  If a sort expression is supplied then it will be used 
to sort the keys matching the filter expression.

_key_ also accepts atone to two additional The "keys" option

## Usage

Three examples of usage are shown below - return all keys 
(unsorted), return all keys sorted by descending `.family_name`, 
return only keys where the `.group` is `"alumni"` sorted 
by ascending `.family_name`.

```shell
    dataset COLLECTION_NAME keys
    dataset COLLECTION_NAME keys true '-.family_name'
    dataset COLLECTION_NAME keys '(eq .group "alumni")' '+.family_name'
```

## filter expressions

A *filter expression* is based on the Go template conditional 
expressions. It uses a prefix notation for the logic (e.g. 
eq - equal, ne - not equal, lt - less than, gt greater than) 
and the value(s) to be compared in [dotpath notation](dotpath.html).

Filters can be simple expressions that result in "true" or 
"false" or compound expressions (e.g. expressions combined with 
_and_ and _or_) that evaluate to "true" or "false".  Simple 
expressions can isolated by parenthasis 
(e.g. `(and (eq .i 1) (ne .s "1") (ne .s "one"))`).

Example filter operators

+ eq - equal (must be same type and value, e.g. 1 does not equal "1")
+ nq - not equal (comparing same type but different values)
+ lt - less than
+ gt - greater than
+ match - given a regular expression and string data return true if they match
+ and - allows you to combine two expression and if both true the expression is true.
+ or - allows you to combine two or more expressions where one is true then expression is true.

#### Simple

A field, `.family_name`, matches a known value, "Feynman".

```
	'(eq .family_name "Feynman")'
```

A field, `.family_name`, does not match a known value, "Feynman".

A field, `.family_name`, does not match value

```
	'(ne .family_name "Feynman")'
```

A field, `.family_name`, match the regular expression `Feym*n`.

```
	'(match "Feynm*n" .family_name)'
```


#### Compound

Two fields match, `.family_name` and `.given_name`, known values "Feynman" and "Richard".

```
	'(and (eq .family_name "Feynman") (eq .given_name "Richard"))'
```

NOTE: That the filters experessions are data type aware. So 
"1" is not the same as 1. Likewise 1 is not the same as 1.0.

## sort expressions

The "keys" option provides for simple one level sorting.
Sorting is described by a plus or minus followed by a dotpath 
to a simple field type (i.e. string, int, or float JSON types). 
In our previous examples sorting ascending by `.family_name` would
be expressed as `+.family_name`. To sort by descending `.family_name` 
you would use the expression `-.family_name`.  By default we assume 
an ascending sort so in practice you can omit a leading "+".

In this example we listing last names of "Smith" sorting by ascending 
given name. The collection name is "people.ds".

```
    dataset people.ds keys '(eq "Smith" .family_name)' '.given_name'
```

In this example we list last anes of "Smith" sorted by descending 
given name.


```
    dataset people.ds keys '(eq "Smith" .family_name)' '-.given_name'
```

## Getting a "sample" of keys

The _dataset_ command respects an option named `-sample N` where N 
is the size (number) of the keys to include in the sample. The sample 
is taken after any filters are applied but may be less than requested 
size if the the filtered results are few than the sample size.  The 
basic process is to get a set of keys, randomly sort the keys, then 
return the top N number of those keys.


Related topics: [count](count.html), [clone](clone), [clone-sample](clone-sample.html), [frame](frame.html), [grid](grid.html)


