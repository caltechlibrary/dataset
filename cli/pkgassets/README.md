
# pkgassets

## USAGE

```
    pkgassets [OPTIONS] VARIABLE_NAME DIR_HOLDING_ASSETS [VARIABLE_NAME DIR_HOLDING_ASSETS ...]
```


## DESCRIPTION

_pkgassets_ generates a Go source directory whos file assets are embedded in a `map[string][]byte` variable. 
This is useful where you want to embed web content, template source code, help docs and other assets that 
can be used for default behavior in a Go command line program or service. 

The map content is harvested from directory holding the assets to be embedded. By default the
path key starts with a slash and does not include the hosting directory (e.g. htdocs/index.html 
would become /index.html if htdocs was used to harvest assets). The prefix and suffix on the
key can be modified based on _pkgassets_'s command line options.

## OPTIONS 

	-c	comment file to be included
	-comment	comment file to be included
	-example	display example(s)
	-ext	Only include files with matching extension
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-o	output filename
	-output	output filename
	-p	package name, if missing defauls to lowercase of variable name
	-package	package name, if missing defauls to lowercase of variable name
	-strip-prefix	strip the prefix from the map key
	-strip-suffix	strip the suffix from the map key
	-v	display version
	-version	display version


## EXAMPLES

```
    pkgassets MAP_VARAIBLE_NAME NAME_OF_DIRECTORY_HOLDING_ASSETS
```

This will result in a Go of type map[string][]byte holding the assets discovered by walking the directory
tree provided. The map's key will represent a path (beginning with "/") pointing at the asset ingested.

```shell
    pkgassets DefaultSite htdocs
```

Assuming that _htdocs_ held

+ index.html
+ css/site.css

In this example the htdocs directory will be crawled and all the files found harvested as a an asset. The
path in the map will not include htdocs and would result in a Go source file like

```golang
    package defaultsite

    var DefaultSite = map[string][]byte{
        "/index.html": []byte{}, // ... the contents of index.html would be here ...
        "/css/site.css": []byte{}, // ... the contents of css/site.css would be here ...
    }
```

If a package name is not provided then the package name will a lowercase name of the map variable name (e.g. 
"var DefaultSite" becomes "package defaultsite"). Likewise if a output name is not provided then the file
name will be the name of the package plus the ".go" extension.

Related examples: [help](examples/help.html), [htdocs](examples/htdocs.html).

pkgassets v0.0.6
