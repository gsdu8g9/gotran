package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var (
	TRANSLATE_URL = "http://translate.google.com/translate_a/t"
	FIRST_STRING  = regexp.MustCompile(`\[("(?:[^\\"]|\\.)*"),`)
)

type Translator struct {
	srcLang string
	dstLang string
}

func NewTranslator(srcLang, dstLang string) *Translator {
	return &Translator{
		srcLang: srcLang,
		dstLang: dstLang,
	}
}

func (t *Translator) fetchTranslated(src []byte) ([]byte, error) {
	res, err := http.PostForm(TRANSLATE_URL, url.Values{
		"sl":     {t.srcLang},
		"tl":     {t.dstLang},
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
	b, err := t.fetchTranslated(src)
	if err != nil {
		return nil, err
	}
	return t.extractText(b)
}
