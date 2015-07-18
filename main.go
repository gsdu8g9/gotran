package main

import (
	"flag"
	"fmt"
	"github.com/yuya-takeyama/argf"
	"io"
	"io/ioutil"
	"os"
)

func guideToHelp() {
	os.Stderr.WriteString(`
Try 'gotran --help' for more information.
`[1:])
}

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

type Option struct {
	IsHelp    bool
	IsVersion bool
	From      string
	To        string
	Files     []string
}

func ParseOption(args []string) (opt *Option, err error) {
	opt = &Option{}
	f := flag.NewFlagSet("gotran", flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	f.BoolVar(&opt.IsHelp, "h", false, "")
	f.BoolVar(&opt.IsHelp, "help", false, "")
	f.BoolVar(&opt.IsVersion, "v", false, "")
	f.BoolVar(&opt.IsVersion, "version", false, "")

	if err = f.Parse(args); err != nil {
		return nil, err
	}
	switch flag.NArg() {
	case 0:
		return nil, fmt.Errorf("no specify FROM and TO language")
	case 1:
		return nil, fmt.Errorf("no specify TO language")
	}
	opt.From, opt.To = flag.Arg(0), flag.Arg(1)
	opt.Files = flag.Args()[2:]

	return opt, nil
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
	opt, err := ParseOption(os.Args[1:])
	if err != nil {
		return err
	}
	switch {
	case opt.IsHelp:
		usage()
		return nil
	case opt.IsVersion:
		version()
		return nil
	}

	t := NewTranslator(opt.From, opt.To)
	r, err := argf.From(opt.Files)
	if err != nil {
		return err
	}
	return do(t, r)
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, "gotran:", err)
		os.Exit(1)
	}
}
