package main

import (
	"strconv"
	"testing"
)

func TestHasLongTLD(t *testing.T) {
	cases := []struct {
		domain  string
		hasLong bool
	}{
		{"", false},
		{"abc.com", false},
		{"sub.abc.com", false},
		{"*.abc.com", false},
		{"abc.co.uk", true},
		{"sub.abc.co.uk", true},
		{"*.abc.co.uk", true},
		{"abc-more.or.us", true},
	}
	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			got := hasLongTLD(tc.domain)
			if got != tc.hasLong {
				t.Errorf("expected %v but got %v", tc.hasLong, got)
			}
		})
	}
}

func TestPrepareDomainNames(t *testing.T) {
	cases := []struct {
		in, out []string
	}{
		// nil and empty slices acceptable
		{nil, nil},
		{[]string{}, []string{}},

		// simple
		{[]string{"abc.com"}, []string{"abc.com", "www.abc.com"}},
		{[]string{"abc.com", "abc-de.co"}, []string{"abc.com", "www.abc.com", "abc-de.co", "www.abc-de.co"}},

		// www variant after main
		{[]string{"abc.com", "www.abc.com"}, []string{"abc.com", "www.abc.com"}},
		{[]string{"abc.co.jp", "www.abc.co.jp"}, []string{"abc.co.jp", "www.abc.co.jp"}},

		// www variant before main
		{[]string{"www.abc.com", "abc.com"}, []string{"www.abc.com", "abc.com"}},
		{[]string{"www.abc.co.jp", "abc.co.jp"}, []string{"www.abc.co.jp", "abc.co.jp"}},

		// sub-domains
		{[]string{"abc.co.uk"}, []string{"abc.co.uk", "www.abc.co.uk"}},
		{[]string{"domain.k12.or.us"}, []string{"domain.k12.or.us"}},
		{[]string{"sub-domain.domain.k12.or.us"}, []string{"sub-domain.domain.k12.or.us"}},

		// wildcards
		{[]string{"abc.com", "*.abc.com"}, []string{"abc.com", "*.abc.com"}},
		{[]string{"abc.com", "*.abc.com", "www.abc.com"}, []string{"abc.com", "*.abc.com"}},
		{[]string{"abc.co.jp", "*.abc.co.jp"}, []string{"abc.co.jp", "*.abc.co.jp"}},
		{[]string{"abc.io", "*.abc.io", "www.abc.io", "no_dot"}, []string{"abc.io", "*.abc.io", "no_dot", "www.no_dot"}},
	}
	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			got := prepareDomainNames(tc.in)
			if len(got) != len(tc.out) {
				t.Logf("got:  %v", got)
				t.Logf("want: %v", tc.out)
				t.Fatalf("got a length of %d in the slice but expected length %d", len(got), len(tc.out))
			}
			for indx, str := range got {
				if tc.out[indx] != str {
					t.Errorf("at index %d expected %q but got %q", indx, tc.out[indx], str)
				}
			}
		})
	}
}

func TestStringSliceContains(t *testing.T) {
	has := map[string][]string{
		"a": {"a", "b"},
		"e": {"c", "e", "b"},
		"":  {"stuff", ""},
	}
	hasNot := map[string][]string{
		"abc": nil,
		"hl":  {},
		"y":   {"a", "b"},
		"i":   {"c", "e", "b"},
	}
	for str, slice := range has {
		if !stringSliceContains(slice, str) {
			t.Errorf("saying slice %v does not have string %q", slice, str)
		}
	}
	for str, slice := range hasNot {
		if stringSliceContains(slice, str) {
			t.Errorf("saying slice %v has string %q", slice, str)
		}
	}
}

func TestRemoveStringsSuffixed(t *testing.T) {
	cases := []struct {
		slice  []string
		suffix string
		output []string
	}{
		{nil, "anything", nil},
		{[]string{"s.abc", "*.abc", "def"}, ".abc", []string{"*.abc", "def"}},
		{[]string{"abc.com", "*.abc.com", "d.co"}, ".abc.com", []string{"abc.com", "*.abc.com", "d.co"}},
		{[]string{"abc.io", "*.gh.ca.us", "www.gh.ca.us"}, ".gh.ca.us", []string{"abc.io", "*.gh.ca.us"}},
		{[]string{"abc", "def", "ghi"}, "-", []string{"abc", "def", "ghi"}},
	}
	for i, tc := range cases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			out := removeStringsSuffixed(tc.slice, tc.suffix)
			if len(out) != len(tc.output) {
				t.Fatalf("got unequal lengths; expected %v but got %v", tc.output, out)
			}
			for j := range out {
				if out[j] != tc.output[j] {
					t.Errorf("got %q at %v but got %q", out[j], tc.output[j], j)
				}
			}
		})
	}
}
