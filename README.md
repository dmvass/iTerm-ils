# iTerm ils

`ils` is a simple implementation of UNIX basic command `ls` with file icons.
The main goal to create this tool is to facilitate research of directories with
set of files with different extensions and ensure similar behavior as `ls`
command.

This tool working only in iTerm2 since it's able to display images within the
terminal. You can use it in other terminal emulators too but without file icons.
This tool working much more slower than standard `ls` implementation but it's
only an iTerm image protocol feature.

![ils bare](/images/bare.png?raw=true)

![ils long](/images/long.png?raw=true)

## More about Iterm File Transfers
By omitting the inline argument (or setting its value to 0), files will be
downloaded and saved in the Downloads folder instead of being displayed inline.
Any kind of file may be downloaded, but only images will display inline. Any
image format that macOS supports will display inline, including PDF, PICT, EPS,
or any number of bitmap data formats (PNG, GIF, etc.). A new menu item titled
Downloads will be added to the menu bar after a download begins, where progress
can be monitored and the file can be located, opened, or removed.

## Instalation
```bash
go get -u github.com/kandziu/iTerm-ils
```

You can easy build binary file:
```bash
make build
```

If you want to build binary file with your custom theme you can use `ILS_THEME_PATH`
environment variable
```bash
ILS_THEME_PATH=path/to/your/theme... make build
```

## How it works
Unix and Unix-like operating systems maintain the idea of a current working
directory, that is, where one is currently positioned in the hierarchy of
directories. When invoked without any arguments, `ils` lists the files in the
current working directory. If another directory is specified, then `ils` will
list the files there, and in fact the user may specify any list of files and
directories to be listed.

Files whose names start with `.` are not listed, unless the `-a` flag is specified,
the `-A` flag is specified, or the files are specified explicitly.

Without options, `ils` displays files in a bare format. This bare format however
makes it difficult to establish the type, permissions, and size of the files.
The most common options to reveal this information or change the list of files are:

|Option    |Description    |
|----------|---------------|
| -l |displaying Unix file types, permissions, number of hard links, owner, group, size, last-modified date and filename|
| -f |do not sort. Useful for directories containing large numbers of files|
| -F |appends a character revealing the nature of a file, for example, * for an executable, or / for a directory|
| -a |lists all files in the given directory, including those whose names start with "." (which are hidden files in Unix)|
| -R |recursively lists subdirectories. The command ls -R / would therefore list all files|
| -t |sort the list of files by modification time|
| -h |print sizes in human readable format|

### Example commands:
```bash
$ ils -l
drwxr--r--   1 fred  editors   4096  drafts
-rw-r--r--   1 fred  editors  30405  edition-32
-r-xr-xr-x   1 fred  fred      8460  edit

$ ils -F
drafts/
edition-32
edit*
```

## Theme
Default theme of this tool was inspired by the my favorite VSCODE file icons
extension [vscode-icons](https://github.com/vscode-icons/vscode-icons). 

## Changes
A full changelog is maintained in the [CAHNGELOG](https://github.com/kandziu/iTerm-ils/blob/master/CHANGELOG.md) file.

## Contributing 
**elasticsearch-partition** is an open source project and contributions are
welcome! Check out the [Issues](https://github.com/kandziu/iTerm-ils/issues)
page to see if your idea for a contribution has already been mentioned, and feel
free to raise an issue or submit a pull request.