
# Creating help docs

Here is the basic recipe for turning text or markdown documents into `map[string][]byte` object
that can drive your help topics.

1. mkdir a docs directory
2. create markdown (files ending in ".md") documents in that directory for your help topics
    + filename without extension will become keyword used with the help command
3. run _pkgassets_ over the directory generating an *asset.go* file in the same folder as your cli program
4. compile your cli progam

## Example

In this example the *cmd/helloworld/helloworld.go* would contains the "main" package for a cli program
you're going to build.  The documentation for _helloworld_ is in a folder called docs.

```
    pkgassets -o cmd/helloworld/assets.go -p main -ext=".md" -strip-prefix="/" -strip-suffix=".md" Help docs
```

This will create the _assets.go_ file which will contain a map[string][]byte of your help docs.
You can then compile _helloworld_ normaly with _go_. Note the pkgassets strips the "docs" from
the value passed in as the key bit not the "/". This is to support using `map[string][]byte`
as holders of web content. We use the additional option "-strip-prefix" to remove the leading
slash leaving the renaming filename as the key in for the mapped help page. Likewise if we have
other documents in the *docs* directory tree we can restrict the help documents harvested to 
a single file type by file extension (e.g. -ext=".md" restricts to markdown files using the ".md"
file extension).

