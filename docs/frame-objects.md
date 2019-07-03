
# frame-objects

## Usage

```
    frame-objects COLLECTION FRAME_NAME
```

Returns the object list of a frame.

## OPTIONS

-p, -pretty
: pretty print JSON output

## Example

If I want to get a list of objects (JSON array of objects) 
for a frame named "captions-dates-locations" from my collection
called "photos.ds" I would do the following (will be using the
`-p` option to pretty print the results)

```
    dataset frame-objects -p photos.ds captions-dates-locations
```



