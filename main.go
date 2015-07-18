package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/yuya-takeyama/argf"
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
	-e, --expr=TEXT  translate text
	-h, --help       show this help message
	-v, --version    print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.3
`[1:])
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "gotran:", err)
}

type Option struct {
	Expr      string `short:"e" long:"expr"`
	IsHelp    bool   `short:"h" long:"help"`
	IsVersion bool   `short:"v" long:"version"`
	From      string
	To        string
	Reader    io.Reader
}

func ParseOption(args []string) (opt *Option, err error) {
	opt = &Option{}
	f := flags.NewParser(opt, flags.PassDoubleDash)

	leave, err := f.ParseArgs(args)
	if err != nil {
		return nil, err
	}
	switch len(leave) {
	case 0:
		return nil, fmt.Errorf("no specify FROM and TO language")
	case 1:
		return nil, fmt.Errorf("no specify TO language")
	}

	opt.From, opt.To = leave[0], leave[1]
	if opt.Expr != "" {
		expr, err := strconv.Unquote(`"` + opt.Expr + `"`)
		if err != nil {
			return nil, err
		}
		opt.Reader = strings.NewReader(expr)
		return opt, nil
	}
	opt.Reader, err = argf.From(leave[2:])
	if err != nil {
		return nil, err
	}
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

func _main() int {
	opt, err := ParseOption(os.Args[1:])
	if err != nil {
		printErr(err)
		guideToHelp()
		return 2
	}
	switch {
	case opt.IsHelp:
		usage()
		return 0
	case opt.IsVersion:
		version()
		return 0
	}

	t := NewTranslator(opt.From, opt.To)
	if err = do(t, opt.Reader); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
