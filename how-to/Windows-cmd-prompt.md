Windows 11 Notes
================

Windows cmd (command prompt) presents some challenges for working with JSON on the command line. This is particularly true on how single and double quotes are handled. Here's an example works under Unix, macOS but not cmd.

```shell
    dataset create T1.ds one '{"one":1}'
```

To get the command to work under cmd you have to type it in like this.

```cmd
    dataset create T1.ds one "{"""one""":1}
```

Fortunately there are alternatives on Windows. PowerShell is what you want to use instead of trying to always sort out strange quote behavior. By Windows 10 it quite stable and is now even cross platform (meaning you can take your PowerShell knowledge to macOS and Linux if you like). With Powershell the orignal example works find.

```pwsh
    dataset create T1.ds one '{"one":1}'
```

There are quirks still lurking. Windows does not ship with the Unix `cat` command. Fortunately PowerShell has one built in and it works regardless if you're running PowerShell on Windows or anther operating system.

Here's what would normally do on a Unix system. In this example the file `one.json` holds our JSON object we want to save into our collection.

~~~shell
cat one.json | dataset create T1.ds one
~~~

In PowerShell I would do this

~~~pwsh
Get-Content one.json | dataset create T1.ds one
~~~
