
# Turning an htdocs directory into an asset

Say you want to create a standalone binary containing the contents of an htdocs directory. The program
name will be called _helloserver_ with its main defined in *cmd/helloserver/helloserver.go*.

Building the package for the all the directory (including sub directories of htdocs can be done
like

```
    pkgassets -o cmd/helloserver/assets.go -p main Htdocs htdocs
```

In the file *helloserver.go* you'd reference the contains of the package variable as `main.Htdocs` passing in the
path requested by your server

```
    func HandleHtdoc(res http.ResponseWriter, req *http.Request) {
        p := req.URL.Path
        if buf, ok := main.Htdocs[p]; ok == true {
            // Probable want to set Content-Type, etc before handing back the data
            io.Write(res, buf)
        } else {
            // Handle your error here
        }
    }
```

Notice the _pkgassets_ by default strips the initial directory name from the path of the value stored. This is
so the path matches easily what is passed in via `req.URL.Path`. Additionally we're not restricting the harvest
to a specific file type like we did in the help example.

