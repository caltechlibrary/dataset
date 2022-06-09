dsv1
====

This module provides limited read only support for v1 based dataset
collections. It is intended to enable migration of data from
a version 1 dataset collection to the current version of the dataset
collection.

Notes
-----

Version 1 dataset collections are pairtree based. The JSON objects
are stored in the pairtree and the attached documents under that
directory.  The operational metadata as well as the general metadata
is maintained in the collection.json file in the root collection
folder. Frame data is also in the root folder.  Frames do not need
to migrate, they tend to be operationally empheral in practice. The
JSON documents and their attachments do need to migrate.

If the version 1 collection is in good working order then the
collections.json file can be used to remap the old collection to
the new one. On the otherhand if it is not in good working order the
pairtree itself can be used to derive the keys and for the objects
in the collection. Leverage the pairtree itself is safeist way to
proceed.


