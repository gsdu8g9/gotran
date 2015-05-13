gotran
===
Translate FILE(s), or standard input.

```
$ echo Hello | gotran en ja
こんにちは
```

Installation
-----
`gotran` can be easily installed as an executable.
Download the latest
[compiled binaries](https://github.com/kusabashira/gotran/releases)
and put it anywhere in your executable path.

Or, if you've done Go development before
and your $GOPATH/bin directory is already in your PATH:
```
$ go get github.com/kusabashira/gotran
```

Usage
------
```
$ gotran [OPTION]... FROM TO [FILE]...

Options:
	--help       show this help message
	--version    print the version
```

Example
------

```
$ echo -e "こんにちは\n世界" | gotran ja en
Hello
World
```

```
$ cat foo
Hello
$ cat bar
World
$ gotran en ja foo bar
こんにちは
世界
```

License
--------
MIT License

Author
-------
wara <kusabashira227@gmail.com>