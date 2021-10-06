Release 1.0.2:

Added support for __datasetd__, a localhost web service for
dataset collections.

Migrated cli package into dataset repository sub-package "github.com/caltechlibrary/dataset/cli". Eventually this package will be replaced by "datasetCli.go" in the root folder.

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

