gotran
======

[![Build Status](https://travis-ci.org/kusabashira/gotran.svg?branch=master)](https://travis-ci.org/kusabashira/gotran)

Translate FILE(s), or standard input.

```
$ echo Hello | gotran en ja
こんにちは
```

It use same way as
[pawurb/termit](https://github.com/pawurb/termit)
to fetch translated result.

Installation
------------

###compiled binaries

See [releases](https://github.com/kusabashira/gotran/releases)

###go get

	$ go get github.com/kusabashira/gotran

Usage
-----

```
$ gotran [OPTION]... FROM TO [FILE]...

Options:
  -e, --expr=TEXT           translate TEXT
  -h, --help                show this help message
  -v, --version             print the version
```

Language
--------

- english - en
- japanese - ja
- polish - pl
- french - fr
- spanish - es
- slovakian - sk
- chinese - zh
- russian - ru
- automatic source language detection - auto

Other language is
[here](https://developers.google.com/translate/v2/using_rest#language-params).

Example
-------

```
$ echo -e "こんにちは\n世界" | gotran ja en
Hello
World

$ gotran ja en -e "こんにちは\n世界"
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
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
