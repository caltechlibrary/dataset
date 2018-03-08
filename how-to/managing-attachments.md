
## Managing File Attachments

In the following examples we are using a collection named *characters.ds*.

Adding high-capt-jack.txt as an attachment to "capt-jack"

```shell
   echo "Hi Capt. Jack, Hello World!" > high-capt-jack.txt
   dataset characters.ds attach capt-jack high-capt-jack.txt
```

List attachments for "capt-jack"

```shell
   dataset characters.ds attachments capt-jack
```

Detach all attachments for "capt-jack" (this will untar in your current directory)

```shell
   dataset characters.ds detach capt-jack
```

Prune high-capt-jack.txt from "capt-jack"

```shell
    dataset characters.ds prune capt-jack high-capt-jack.txt
```

Prune all attachments from "capt-jack"

```shell
   dataset prune capt-jack
```

