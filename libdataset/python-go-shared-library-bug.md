
# Go based C-Shared libraries on Mac OS X 10.x

When I use multiple Go based C-Shared libraries from Python 3.6 I get a runtime error, "runtime/cgo: could not obtain pthread_keys". 
This appears to be a Go linking issue for Mac OS X 10.x and Go 1.10.x. It appears that when the go runtime starts on the second 
shared library that it can't find the **pthread_keys** entry point.  This appears to be a known but unresolved, see issue
[#17200](https://github.com/golang/go/issues/17200) in Go's Github repository. It appears this problem goes back several years
and re-occurs in numerious forms.  It doesn't appear to be solvabled without seriously hacking the Go linker based on the version 
of Go an OS you're usiing.


## Additional references

+ https://go-review.googlesource.com/c/go/+/108679
+ http://grokbase.com/t/gg/golang-nuts/159wwdy87y/go-nuts-how-can-i-do-to-fix-%60runtime-cgo-could-not-obtain-pthread-keys-on-darwin-amd64

