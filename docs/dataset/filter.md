
# filter

## Description

filter is based on the Go template conditional expressions. It uses a prefix
notation for the logic (e.g. eq - equal, ne - not equal, lt - less than, gt 
greater than) and the value(s) to be compared in dotpath notation.

Filters can be simple expressions that result in "true" or "false" or compound
expressions (e.g. expressions combined with _and_ and _or_) 
that evaluate to "true" or "false".  Simple expressions can isolated by 
parenthasis (e.g. `(and (eq .i 1) (ne .s "1") (ne .s "one"))`).


### NOTE 

That the filters experessions are data type aware. So "1" is not the same
as 1. Likewise 1 is not the same as 1.0.

## Expresions

Expressions that make up a filter use prefix notation (e.g. the operator
is first followed by values).  They can be be simple or compound.

+ eq - equal (must be same type and value, e.g. 1 does not equal "1")
+ nq - not equal (comparing same type but different values)
+ lt - less than
+ gt - greater than
+ and - allows you to combine two expression and if both true the expression is true.
+ or - allows you to combine two or more expressions where one is true then expression is true.

### Simple

A field, .family_name, matches a known value, "Feynman".


```
	'(eq .family_name "Feynman")'
```

A field, .family_name, does not match a known value, "Feynman".

A field, .family_name, does not match value

```
	'(ne .family_name "Feynman")'
```

### Compound

Two fields match, .family_name and .given_name, known values "Feynman" and "Richard".

```
	'(and (eq .family_name "Feynman") (eq .given_name "Richard"))'
```

Related topics: extract and export

