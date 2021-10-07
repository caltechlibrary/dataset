Release 1.0.2:

Added support for __datasetd__, a localhost web service for
dataset collections. The web service supports a subset of
the command line tool.

Both __datasetd__ and __dataset__ command line program now
include a "lock.pid" file in the collection root. This is to
prevent multiple processes from clashing when maintaining the
"collections.json" file in the collection root.

Migrated cli package into dataset repository sub-package "github.com/caltechlibrary/dataset/cli". Eventually this package will be replaced by "datasetCli.go" in the root folder.

In the dataset command line program the verb "detach" has been
renamed "retrieve" better describe the action. "detach" is depreciated
and will be removed in upcoming releases.

Release 1.0.1:

- Keys are stored lowercase
- Removed filtering and sorting options from dataset and libdataset
- Use pairtree 1.0.2 configurable separator
- Added check and repair for migrating to case insensitive keys and path
- Updated required packages to latest releases
- Added notes about Windows cmd prompt issues when providing JSON objects on command line
- Added M1 support for libdataset

Release 1.0.0:

- Initial Stable Release

