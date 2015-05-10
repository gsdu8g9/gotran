package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
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
v0.1.0
`[1:])
}

var (
	TRANSLATE_URL = "http://translate.google.com/translate_a/t"
	FIRST_ARRAY   = regexp.MustCompile(`^\[\[\[.+?\]\],,"`)
	FIRST_STRING  = regexp.MustCompile(`\["((?:[^\\"]|\\.)*)",`)
	NEW_LINE      = regexp.MustCompile(`\\n`)
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
		"sl":     []string{t.from},
		"tl":     []string{t.to},
		"ie":     []string{"UTF-8"},
		"oe":     []string{"UTF-8"},
		"client": []string{"t"},
		"text":   []string{string(src)},
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

func (t *Translator) extractText(b []byte) []byte {
	var buf [][]byte
	for _, a := range FIRST_ARRAY.FindAllSubmatch(b, -1) {
		for _, s := range FIRST_STRING.FindAllSubmatch(a[0], -1) {
			s[1] = NEW_LINE.ReplaceAll(s[1], []byte("\n"))
			buf = append(buf, s[1])
		}
	}
	return bytes.Join(buf, []byte(""))
}

func (t *Translator) Translate(src []byte) ([]byte, error) {
	b, err := t.fetchResult(src)
	if err != nil {
		return nil, err
	}
	return t.extractText(b), nil
}

func do(t *Translator, f *os.File) error {
	src, err := ioutil.ReadAll(f)
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

func _main() (int, error) {
	isHelp := flag.Bool("help", false, "")
	isVersion := flag.Bool("version", false, "")
	flag.Usage = usage
	flag.Parse()
	switch {
	case *isHelp:
		usage()
		return 2, nil
	case *isVersion:
		version()
		return 2, nil
	case flag.NArg() < 1:
		return 1, fmt.Errorf("no specify FROM and TO language")
	case flag.NArg() < 2:
		return 1, fmt.Errorf("no specify TO language")
	}

	t := NewTranslator(flag.Arg(0), flag.Arg(1))
	if flag.NArg() < 3 {
		if err := do(t, os.Stdin); err != nil {
			return 1, err
		}
		return 0, nil
	}
	for _, name := range flag.Args()[2:] {
		f, err := os.Open(name)
		if err != nil {
			return 1, err
		}
		defer f.Close()
		if err = do(t, f); err != nil {
			return 1, err
		}
	}
	return 0, nil
}

func main() {
	exitCode, err := _main()
	if err != nil {
		fmt.Fprintln(os.Stderr, "gotran:", err)
	}
	os.Exit(exitCode)
}
