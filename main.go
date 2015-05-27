package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

func usage() {
	os.Stderr.WriteString(`
Usage: gotran [OPTION]... FROM TO [FILE]...
Translate FILE(s), or standard input.

Options:
	--help       show this help message
	--version    print the version
`[1:])
}

func version() {
	os.Stderr.WriteString(`
v0.1.3
`[1:])
}

var (
	TRANSLATE_URL = "http://translate.google.com/translate_a/t"
	FIRST_STRING  = regexp.MustCompile(`\[("(?:[^\\"]|\\.)*"),`)
)

type Translator struct {
	from string
	to   string
}

func NewTranslator(from, to string) *Translator {
	return &Translator{
		from: from,
		to:   to,
	}
}

func (t *Translator) fetchResult(src []byte) ([]byte, error) {
	res, err := http.PostForm(TRANSLATE_URL, url.Values{
		"sl":     {t.from},
		"tl":     {t.to},
		"ie":     {"UTF-8"},
		"oe":     {"UTF-8"},
		"client": {"t"},
		"text":   {string(src)},
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (t *Translator) extractText(b []byte) ([]byte, error) {
	var buf [][]byte

	a := b[:bytes.Index(b, []byte("]],"))]
	for _, s := range FIRST_STRING.FindAllSubmatch(a, -1) {
		t, err := strconv.Unquote(string(s[1]))
		if err != nil {
			return nil, err
		}
		buf = append(buf, []byte(t))
	}
	return bytes.Join(buf, []byte("")), nil
}

func (t *Translator) Translate(src []byte) ([]byte, error) {
	b, err := t.fetchResult(src)
	if err != nil {
		return nil, err
	}
	return t.extractText(b)
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
	isHelp := flag.Bool("help", false, "")
	isVersion := flag.Bool("version", false, "")
	flag.Usage = usage
	flag.Parse()
	switch {
	case *isHelp:
		usage()
		return nil
	case *isVersion:
		version()
		return nil
	case flag.NArg() < 1:
		return fmt.Errorf("no specify FROM and TO language")
	case flag.NArg() < 2:
		return fmt.Errorf("no specify TO language")
	}

	t := NewTranslator(flag.Arg(0), flag.Arg(1))
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
