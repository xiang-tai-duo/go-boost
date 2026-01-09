// Package boost
// File:        string.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/string.go
// Author:      Vibe Coding
// Created:     12/30/2025 11:03:46
// Description: STRING is a wrapper for string operations, providing a set of methods for string manipulation.
// --------------------------------------------------------------------------------
package boost

import (
	"strings"
	"unicode"
)

type (
	STRING struct {
		value string
	}
)

func NewString(value string) *STRING {
	return &STRING{value: value}
}

func (s *STRING) Equals(t string) bool {
	return s.Compare(t) == 0
}

func (s *STRING) EqualsIgnoreCase(t string) bool {
	return s.EqualFold(t)
}

func (s *STRING) Compare(t string) int {
	return strings.Compare(s.value, t)
}

func (s *STRING) Contains(substr string) bool {
	return strings.Contains(s.value, substr)
}

func (s *STRING) ContainsAny(chars string) bool {
	return strings.ContainsAny(s.value, chars)
}

func (s *STRING) ContainsRune(r rune) bool {
	return strings.ContainsRune(s.value, r)
}

func (s *STRING) Count(substr string) int {
	return strings.Count(s.value, substr)
}

func (s *STRING) Cut(sep string) (before STRING, after STRING, found bool) {
	b, a, f := strings.Cut(s.value, sep)
	return STRING{value: b}, STRING{value: a}, f
}

func (s *STRING) CutPrefix(prefix string) (STRING, bool) {
	result, found := strings.CutPrefix(s.value, prefix)
	return STRING{value: result}, found
}

func (s *STRING) CutSuffix(suffix string) (STRING, bool) {
	result, found := strings.CutSuffix(s.value, suffix)
	return STRING{value: result}, found
}

func (s *STRING) EqualFold(t string) bool {
	return strings.EqualFold(s.value, t)
}

func (s *STRING) Fields() []string {
	return strings.Fields(s.value)
}

func (s *STRING) FieldsFunc(f func(rune) bool) []string {
	return strings.FieldsFunc(s.value, f)
}

func (s *STRING) HasPrefix(prefix string) bool {
	return strings.HasPrefix(s.value, prefix)
}

func (s *STRING) HasSuffix(suffix string) bool {
	return strings.HasSuffix(s.value, suffix)
}

func (s *STRING) Index(substr string) int {
	return strings.Index(s.value, substr)
}

func (s *STRING) IndexAny(chars string) int {
	return strings.IndexAny(s.value, chars)
}

func (s *STRING) IndexByte(c byte) int {
	return strings.IndexByte(s.value, c)
}

func (s *STRING) IndexFunc(f func(rune) bool) int {
	return strings.IndexFunc(s.value, f)
}

func (s *STRING) IndexRune(r rune) int {
	return strings.IndexRune(s.value, r)
}

func (s *STRING) Join(elems []string) STRING {
	return STRING{value: strings.Join(elems, s.value)}
}

func (s *STRING) LastIndex(substr string) int {
	return strings.LastIndex(s.value, substr)
}

func (s *STRING) LastIndexAny(chars string) int {
	return strings.LastIndexAny(s.value, chars)
}

func (s *STRING) LastIndexByte(c byte) int {
	return strings.LastIndexByte(s.value, c)
}

func (s *STRING) LastIndexFunc(f func(rune) bool) int {
	return strings.LastIndexFunc(s.value, f)
}

func (s *STRING) Left(length int) STRING {
	var result STRING
	if length >= len(s.value) {
		result = STRING{value: s.value}
	} else {
		result = STRING{value: s.value[:length]}
	}
	return result
}

func (s *STRING) Map(mapping func(rune) rune) STRING {
	return STRING{value: strings.Map(mapping, s.value)}
}

func (s *STRING) Mid(start int, length int) STRING {
	var result STRING
	if start < 0 {
		start = 0
	}
	if start >= len(s.value) {
		result = STRING{value: ""}
	} else {
		end := start + length
		if end > len(s.value) {
			end = len(s.value)
		}
		result = STRING{value: s.value[start:end]}
	}
	return result
}

func (s *STRING) Repeat(count int) STRING {
	return STRING{value: strings.Repeat(s.value, count)}
}

func (s *STRING) Replace(old string, newStr string, n int) STRING {
	return STRING{value: strings.Replace(s.value, old, newStr, n)}
}

func (s *STRING) ReplaceAll(old string, newStr string) STRING {
	return STRING{value: strings.ReplaceAll(s.value, old, newStr)}
}

func (s *STRING) Right(length int) STRING {
	var result STRING
	if length >= len(s.value) {
		result = STRING{value: s.value}
	} else {
		result = STRING{value: s.value[len(s.value)-length:]}
	}
	return result
}

func (s *STRING) Split(sep string) []string {
	return strings.Split(s.value, sep)
}

func (s *STRING) SplitAfter(sep string) []string {
	return strings.SplitAfter(s.value, sep)
}

func (s *STRING) SplitAfterN(sep string, n int) []string {
	return strings.SplitAfterN(s.value, sep, n)
}

func (s *STRING) SplitN(sep string, n int) []string {
	return strings.SplitN(s.value, sep, n)
}

func (s *STRING) String(value ...string) string {
	if len(value) > 0 {
		s.value = value[0]
	}
	return s.value
}

func (s *STRING) SubStr(start int) STRING {
	var result STRING
	if start < 0 {
		start = 0
	}
	if start >= len(s.value) {
		result = STRING{value: ""}
	} else {
		result = STRING{value: s.value[start:]}
	}
	return result
}

func (s *STRING) Title() STRING {
	return STRING{value: strings.ToTitle(s.value)}
}

func (s *STRING) ToLower() STRING {
	return STRING{value: strings.ToLower(s.value)}
}

func (s *STRING) ToLowerSpecial(c unicode.SpecialCase) STRING {
	return STRING{value: strings.ToLowerSpecial(c, s.value)}
}

func (s *STRING) ToTitleSpecial(c unicode.SpecialCase) STRING {
	return STRING{value: strings.ToTitleSpecial(c, s.value)}
}

func (s *STRING) ToUpper() STRING {
	return STRING{value: strings.ToUpper(s.value)}
}

func (s *STRING) ToUpperSpecial(c unicode.SpecialCase) STRING {
	return STRING{value: strings.ToUpperSpecial(c, s.value)}
}

func (s *STRING) Trim(cutset string) STRING {
	return STRING{value: strings.Trim(s.value, cutset)}
}

func (s *STRING) TrimFunc(f func(rune) bool) STRING {
	return STRING{value: strings.TrimFunc(s.value, f)}
}

func (s *STRING) TrimLeft(cutset string) STRING {
	return STRING{value: strings.TrimLeft(s.value, cutset)}
}

func (s *STRING) TrimLeftFunc(f func(rune) bool) STRING {
	return STRING{value: strings.TrimLeftFunc(s.value, f)}
}

func (s *STRING) TrimPrefix(prefix string) STRING {
	return STRING{value: strings.TrimPrefix(s.value, prefix)}
}

func (s *STRING) TrimRight(cutset string) STRING {
	return STRING{value: strings.TrimRight(s.value, cutset)}
}

func (s *STRING) TrimRightFunc(f func(rune) bool) STRING {
	return STRING{value: strings.TrimRightFunc(s.value, f)}
}

func (s *STRING) TrimSpace() STRING {
	return STRING{value: strings.TrimSpace(s.value)}
}

func (s *STRING) TrimSuffix(suffix string) STRING {
	return STRING{value: strings.TrimSuffix(s.value, suffix)}
}
