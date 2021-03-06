package vslparser

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type kv struct {
	k string
	v string
}

// TestSplitLine tests that splitLine correctly handles white-space and missing
// keys/values in the input line and that white-space is preserved at the tail
// of the value.
func TestSplitLine(t *testing.T) {
	samples := map[string]kv {
		"": kv{},
		"    foo    ": kv{k: "foo"},
		"foo bar": kv{k: "foo", v: "bar"},
		"  foo  bar": kv{k: "foo", v: "bar"},
		" foo    bar   ": kv{k: "foo", v: "bar   "},
		"				 foo	bar	 ": kv{k: "foo", v: "bar	 "},
	}
	for line, kv := range samples {
		gk, gv := splitLine(line)
		if gk != kv.k {
			t.Errorf("parsing %q should give %q as key, got %q", line, kv.k, gk)
		}
		if gv != kv.v {
			t.Errorf("parsing %q should give %q as value, got %q", line, kv.v, gv)
		}
	}
}

// stringScanner is a helper which provides suitable bufio.Scanner object that
// can be passed to Parse so that the given string is processed by the parser.
func stringScanner(s string) *bufio.Scanner {
	return bufio.NewScanner(strings.NewReader(s))
}

// testParseOK tests that s is parsed as e without errors.
func testParseOK(t *testing.T, e *Entry, s string) {
	got, err := Parse(stringScanner(s))
	if err != nil {
		t.Errorf("failed to parse %q: %v", s, err)
		return
	}
	if !reflect.DeepEqual(e, got) {
		t.Errorf("parsing %q should give %v, got %v", s, e, got)
	}
}

// testParseMultipleOK tests that s is parsed as a chain of entries from e
// without errors.
func testParseMultipleOK(t *testing.T, e []*Entry, s string) {
	scanner := stringScanner(s)
	for _, ent := range e {
		got, err := Parse(scanner)
		if err != nil {
			t.Errorf("failed to parse %q: %v", s, err)
			return
		}
		if !reflect.DeepEqual(ent, got) {
			t.Errorf("parsing %q should give %v, got %v", s, ent, got)
		}
	}
}

// testParseError tests that parsing s produces an error.
func testParseError(t *testing.T, s string) {
	_, err := Parse(stringScanner(s))
	if err == nil {
		t.Errorf("parsing %q should be a parse error", s)
	} else {
		t.Logf("parsing %q gives: %v", s, err)
	}
}

// TestParse tests that various inputs are either parsed correctly or produce
// errors.
func TestParse(t *testing.T) {
	testParseOK(t, &Entry{
		Kind: BeReq,
		VXID: 123,
		Fields: Fields{},
	}, "* << BeReq >> 123\n- End")

	testParseOK(t, &Entry{
		Kind: Request,
		VXID: 40000000,
		Fields: Fields{
			"Foo": []string{
				"Bar",
				"Baz",
			},
			"Bar": []string{
				"Foo  Bar    Baz	", // Trailing tab valid.
			},
		},
	}, "*   <<  Request >> 40000000\n- Foo Bar\n-Foo Baz\n- Bar     Foo  Bar    Baz	\n- End")
	testParseMultipleOK(t, []*Entry{
		&Entry{
			Kind:   BeReq,
			VXID:   123,
			Fields: Fields{},
		},
		&Entry{
			Kind:   BeReq,
			VXID:   124,
			Fields: Fields{},
		},
	}, "* << BeReq >> 123\n- End\n\n* << BeReq >> 124\n- End")

	testParseError(t, "")
	testParseError(t, "- ")
	testParseError(t, "* << Request >> 1\n - Foo Bar\n- End")
	testParseError(t, "* << Request >> Foo")
	testParseError(t, "* << Request >> 1")
}

func TestEOF(t *testing.T) {
	f, _ := os.Open(os.DevNull)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	_, err := Parse(scanner)
	if err != io.EOF {
		t.Errorf("parsing should result in an EOF error. Got error: '%s'", err)
	} else {
		t.Logf("parsing properly returned EOF")
	}
}
