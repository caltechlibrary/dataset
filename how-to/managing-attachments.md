Managing File Attachments
=========================

In the following examples we are using a collection named *characters.ds*.

Adding high-capt-jack.txt as an attachment to "capt-jack"

```shell
   echo "Hi Capt. Jack, Hello World!" > high-capt-jack.txt
   dataset attach characters.ds capt-jack high-capt-jack.txt
```

List attachments for "capt-jack"

```shell
   dataset attachments characters.ds capt-jack
```

Detach all attachments for "capt-jack" (this will untar in your current directory)

```shell
   dataset detach characters.ds capt-jack
```

Prune high-capt-jack.txt from "capt-jack"

```shell
    dataset prune characters.ds capt-jack high-capt-jack.txt
```

Prune all attachments from "capt-jack"

```shell
   dataset prune capt-jack
```

