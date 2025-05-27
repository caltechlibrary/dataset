---
title: Migrating from v2 to v3
---

# Migrating from v2 to v3

Dataset v3 is not compatible with dataset v2.  To migrate from one version to another you use the dump from a dataset v2 client then use dataset v3's load.

1. Install latest v2
2. Rename dataset cli to datasetV2
3. Install latest v3
4. Dump usng datasetV2
5. Load using dataset

In this example `my_objects.ds` is the dataset we want to migrate from. The new version is called `new_objects.ds`.

~~~shell
curl -L https://caltechlibrary.github.io/dataset/installer.sh 2.2.4
mv $HOME/bin/dataset $HOME/bin/datasetV2
curl -L https://caltechlibrary.github.io/dataset/installer.sh
datasetV2 dump my_objects.ds >objects.jsonl
dataset init new_objects.ds
dataset load new_objects.ds <objects.jsonl
~~~~

