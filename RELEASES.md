
# Release Process and tags

dataset is going through rapid development and changes as we 
evolve the software to meet the needs of the DLD group in 
Caltech Library.  Below is our current policy regarding releases.

## Preleases and production releases

All releases to should a semantic version number (i.e. [semvar]()).
Releases that are experimental
are tagged "pre-release" on Github. Releases intended to be used in our production systems or used by other library staff should not have the "pre-release" option checked on GitHub.

### Pre-releases

Pre-releases may or may not have zip'ed executables ready 
for installion.  Where they do we are currently targetting Linux on AMD64, 
Raspbian on ARM 7, Mac OS X on AMD64 and Windows 10 (for use from 
Window's command prompt) on AMD 64.  From time to time preleases may 
also include an experiemental Python module compiled for Mac OS X and 
Linux.

### Production releases

Production releases will include zip files for installing pre-compiled
binaries for Linux on AMD64, Mac OS X AMD 64, Windows 10 on AMD 64 and
Raspberry Pi on ARM 7. Production release may include experiement code
or utilities like the Python module for dataset.

## Making a release

1. Set the version number in PACKAGE.go (where PACKAGE is the name of the package, e.g. dataset is the name of the dataset
package so you'd change the version number in dataset.go).
2. Run `make clean`
3. Run `make test` and make sure they pass, if some fail document it if you plan to release it (e.g. GSheet integration not tested because...)
4. Run `make release`

You are now ready to go to Github and create a release. If you are uploading compiled versions upload the zip files in the _dist_
folder.

