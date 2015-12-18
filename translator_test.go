package main

import (
	"reflect"
	"testing"
)

type extractTest struct {
	Src []byte
	Dst []byte
}

var indexTestsExtract = []extractTest{
	{
		Src: []byte(`[[["こんにちは","Hello","Kon'nichiwa",""]],,"en",,[["こんにちは",[1],false,false,1000,0,1,0]],[["Hello",1,[["こんにちは",1000,false,false],["ハロー",0,false,false],["のhello",0,false,false],["ようこそ",0,false,false]],[[0,5]],"Hello"]],,,[["en"]],2]`),
		Dst: []byte("こんにちは"),
	}, {
		Src: []byte(`[[["これはペンです。","This is a pen.","Kore wa pendesu.",""]],,"en",,[["これは",[1],false,false,997,0,2,0],["ペンです",[2],false,false,993,2,4,0],["。",[3],false,false,1000,4,5,0]],[["This (c:nsubjmain)",1,[["これは",997,false,false],["これが",0,false,false],["これにより",0,false,false]],[[0,4]],"This is a pen."],["pen is",2,[["ペンです",993,false,false]],[[5,7],[10,13]],""],[".",3,[["。",1000,false,false]],[[13,14]],""]],,,[["en"]],6]`),
		Dst: []byte("これはペンです。"),
	}, {
		Src: []byte(`[[["こんにちは。\n","hi.\n","Kon'nichiwa.",""],["こんにちは。","hi.","Kon'nichiwa.",""]],,"en",,[["こんにちは",[1],false,false,1000,0,1,0],["。",[2],false,false,1000,1,2,0],["n",[3],true,false,0,0,0,0],["こんにちは",[7],false,false,1000,0,1,1],["。",[8],false,false,1000,1,2,1]],[["hi",1,[["こんにちは",1000,false,false],["ハイ",0,false,false],["ハイテク",0,false,false],["のHi",0,false,false],["やあ",0,false,false]],[[0,2]],"hi."],[".",2,[["。",1000,false,false]],[[2,3]],""],["n",3,,[[0,1]],"n"],["hi",7,[["こんにちは",1000,false,false],["ハイ",0,false,false],["ハイテク",0,false,false],["のHi",0,false,false],["やあ",0,false,false]],[[0,2]],"hi."],[".",8,[["。",1000,false,false]],[[2,3]],""]],,,[["ca","en"]],16]`),
		Dst: []byte("こんにちは。\nこんにちは。"),
	},
}

func TestExtract(t *testing.T) {
	tr := NewTranslator("en", "ja")
	for _, test := range indexTestsExtract {
		expect := test.Dst
		actual, err := tr.extractText(test.Src)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("got %q; want %q",
				actual, expect)
		}
	}
}

func TestTranslate(t *testing.T) {
	tr := NewTranslator("en", "ja")
	src := []byte("Hello")

	expect := []byte("こんにちは")
	actual, err := tr.Translate(src)
	if err != nil {
		t.Fatalf("%#v.Translate(%q) returns %v, want nil",
			tr, string(src), err)
	}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("%#v.Translate(%q) got %q: want %q",
			tr, string(src), string(actual), string(expect))
	}
}
