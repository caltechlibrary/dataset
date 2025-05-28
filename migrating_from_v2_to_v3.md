---
title: Migrating from v2 to v3
---

# Migrating from v2 to v3

Dataset v3 is not compatible with dataset v2.  To migrate from one version to another you use the dump from a dataset v2 client then use dataset v3's load.

1. Install last v2 dataset version if not previously installed
3. Install latest v3
4. Dump usng dataset
5. Load using dataset3

In this example `old_objects.ds` is the dataset we want to migrate from. The new version is called `new_objects.ds`.

~~~shell
curl -L https://caltechlibrary.github.io/dataset/installer.sh 2.2.4
curl -L https://caltechlibrary.github.io/dataset/installer.sh
dataset dump old_objects.ds >objects.jsonl
dataset3 init new_objects.ds
dataset3 load new_objects.ds <objects.jsonl
~~~~


