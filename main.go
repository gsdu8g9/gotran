package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func usage() {
	os.Stderr.WriteString(`
Usage: gotran [OPTION]... FROM TO [FILE]...
Translate FILE(s), or standard input.

Options:
	-h, --help       show this help message
	-v, --version    print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.3
`[1:])
}

func do(t *Translator, r io.Reader) error {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	dst, err := t.Translate(src)
	if err != nil {
		return err
	}
	os.Stdout.Write(dst)
	os.Stdout.WriteString("\n")
	return nil
}

func _main() error {
	var isHelp, isVersion bool
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isHelp, "help", false, "")
	flag.BoolVar(&isVersion, "v", false, "")
	flag.BoolVar(&isVersion, "version", false, "")
	flag.Usage = usage
	flag.Parse()
	switch {
	case isHelp:
		usage()
		return nil
	case isVersion:
		version()
		return nil
	case flag.NArg() < 1:
		return fmt.Errorf("no specify FROM and TO language")
	case flag.NArg() < 2:
		return fmt.Errorf("no specify TO language")
	}
	from, to := flag.Arg(0), flag.Arg(1)

	t := NewTranslator(from, to)
	if flag.NArg() < 3 {
		return do(t, os.Stdin)
	}

	var in []io.Reader
	for _, name := range flag.Args()[2:] {
		f, err := os.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		in = append(in, f)
	}
	return do(t, io.MultiReader(in...))
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, "gotran:", err)
		os.Exit(1)
	}
}
