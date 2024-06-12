%dataset(1) user manual | version 2.1.12 cfcaeeb
% R. S. Doiel and Tom Morrell
% 2024-06-12

# Compatibity

As of 2.1 a minimal level of backward dataset v1.1 was added. This
includes support for libdataset, a C-shared library.  The goal in
adding support for v1.1 was to facilate migration in the Caltech
Library feeds project.  Some features in v1.1 are not available and
the features recreated from v1.1 were done so in a way that maintains
the approach of v2. Aside from continuing support for libdataset the
specific v1.1 features will likely be depreciated over time in the
name of keeping the project as simple as possible maintaining focus
on the core benefits of dataset versus other JSON store systems.

## What was left out

The methods related to Namaste data have not be implememented as
v2 of dataset uses a codemeta.json file for collection metadata.
E.g. `Who`, `Where`, `Location`, `Contact`.

The methods `KeySort` and `KeyFilter` are not included.

## Changed behavior

The `DocPath` method returns a full path to the JSON documented stored
in a Pairtree collection. If the collection uses an SQL store then you
will get an empty string and an error indicating that storage type is
not supported by `DocPath`.

The `Keys` method returns a list of keys (slice of strings) and an
error value.

## Method name changes

The following methods were normalized to conform with Go idioms in
the standard library.

- `KeyExists` became `HasKey`
- `FrameExists` became `HasFrame`

The methods for working with versioned content have the order of the
parameters revised, basically semver is now the last paramter.  

## Changes method signatures

Do to the changes in how keys and attachments are handled in v2 you
don't need to "santize" your returned JSON object. Dataset v2 does not
inject a `_Key` or `_Attachments` values in the JSON document stored.
This changes the `Read` method signature for collection objects.

The `Init` method takes a DSN as the second parameter, if the DSN is
an empty string then the collection created will use a Pairtree store.


