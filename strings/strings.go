// Package strings
// File:        strings.go
// Url:         https://github.com/xiang-tai-duo/go-boost/blob/master/strings/strings.go
// Author:      Vibe Coding
// Created:     2025/12/20 12:31:58
// Description: STRINGS is a wrapper for string operations, providing a set of methods for string manipulation.
// --------------------------------------------------------------------------------
package strings

import (
	"strconv"
	__strings "strings"
	"unicode"
)

type (
	STRINGS struct {
		value string
	}
)

//goland:noinspection GoUnusedExportedFunction
func New(value string) *STRINGS {
	return &STRINGS{value: value}
}

//goland:noinspection SpellCheckingInspection
func (s *STRINGS) Atoi() (int, error) {
	result := 0
	err := error(nil)
	sign := 1
	index := 0
	for index < len(s.value) && unicode.IsSpace(rune(s.value[index])) {
		index++
	}
	if index < len(s.value) {
		if s.value[index] == '+' {
			index++
		} else if s.value[index] == '-' {
			sign = -1
			index++
		}
	}
	for index < len(s.value) && unicode.IsDigit(rune(s.value[index])) {
		result = result*10 + int(s.value[index]-'0')
		index++
	}
	if index < len(s.value) && !unicode.IsSpace(rune(s.value[index])) {
		result = 0
		err = &strconv.NumError{Func: "Atoi", Num: s.value, Err: strconv.ErrSyntax}
	} else {
		result = sign * result
		err = nil
	}
	return result, err
}

func (s *STRINGS) Compare(t string) int {
	return __strings.Compare(s.value, t)
}

func (s *STRINGS) Contains(substr string) bool {
	return __strings.Contains(s.value, substr)
}

func (s *STRINGS) ContainsAny(chars string) bool {
	return __strings.ContainsAny(s.value, chars)
}

func (s *STRINGS) ContainsRune(r rune) bool {
	return __strings.ContainsRune(s.value, r)
}

func (s *STRINGS) Count(substr string) int {
	return __strings.Count(s.value, substr)
}

func (s *STRINGS) Cut(sep string) (before string, after string, found bool) {
	beforePart, afterPart, foundFlag := __strings.Cut(s.value, sep)
	return beforePart, afterPart, foundFlag
}

func (s *STRINGS) CutPrefix(prefix string) (string, bool) {
	result, found := __strings.CutPrefix(s.value, prefix)
	return result, found
}

func (s *STRINGS) CutSuffix(suffix string) (string, bool) {
	result, found := __strings.CutSuffix(s.value, suffix)
	return result, found
}

func (s *STRINGS) EqualFold(t string) bool {
	return __strings.EqualFold(s.value, t)
}

func (s *STRINGS) Equals(t string) bool {
	return s.Compare(t) == 0
}

func (s *STRINGS) EqualsIgnoreCase(t string) bool {
	return s.EqualFold(t)
}

func (s *STRINGS) Fields() []string {
	return __strings.Fields(s.value)
}

func (s *STRINGS) FieldsFunc(f func(rune) bool) []string {
	return __strings.FieldsFunc(s.value, f)
}

func (s *STRINGS) HasPrefix(prefix string) bool {
	return __strings.HasPrefix(s.value, prefix)
}

func (s *STRINGS) HasSuffix(suffix string) bool {
	return __strings.HasSuffix(s.value, suffix)
}

func (s *STRINGS) Index(substr string) int {
	return __strings.Index(s.value, substr)
}

func (s *STRINGS) IndexAny(chars string) int {
	return __strings.IndexAny(s.value, chars)
}

func (s *STRINGS) IndexByte(c byte) int {
	return __strings.IndexByte(s.value, c)
}

func (s *STRINGS) IndexFunc(f func(rune) bool) int {
	return __strings.IndexFunc(s.value, f)
}

func (s *STRINGS) IndexRune(r rune) int {
	return __strings.IndexRune(s.value, r)
}

//goland:noinspection SpellCheckingInspection
func (s *STRINGS) Itoa(i int) string {
	return strconv.Itoa(i)
}

func (s *STRINGS) Join(elems []string) string {
	return __strings.Join(elems, s.value)
}

func (s *STRINGS) LastIndex(substr string) int {
	return __strings.LastIndex(s.value, substr)
}

func (s *STRINGS) LastIndexAny(chars string) int {
	return __strings.LastIndexAny(s.value, chars)
}

func (s *STRINGS) LastIndexByte(c byte) int {
	return __strings.LastIndexByte(s.value, c)
}

func (s *STRINGS) LastIndexFunc(f func(rune) bool) int {
	return __strings.LastIndexFunc(s.value, f)
}

func (s *STRINGS) Left(length int) string {
	result := ""
	if length >= len(s.value) {
		result = s.value
	} else {
		result = s.value[:length]
	}
	return result
}

func (s *STRINGS) Lower() string {
	return __strings.ToLower(s.value)
}

func (s *STRINGS) LowerSpecial(c unicode.SpecialCase) string {
	return __strings.ToLowerSpecial(c, s.value)
}

func (s *STRINGS) Map(mapping func(rune) rune) string {
	return __strings.Map(mapping, s.value)
}

func (s *STRINGS) Mid(start int, length int) string {
	result := ""
	if start < 0 {
		start = 0
	}
	if start >= len(s.value) {
		result = ""
	} else {
		end := start + length
		if end > len(s.value) {
			end = len(s.value)
		}
		result = s.value[start:end]
	}
	return result
}

func (s *STRINGS) Repeat(count int) string {
	return __strings.Repeat(s.value, count)
}

func (s *STRINGS) Replace(old string, newStr string, n int) string {
	return __strings.Replace(s.value, old, newStr, n)
}

func (s *STRINGS) ReplaceAll(old string, newStr string) string {
	return __strings.ReplaceAll(s.value, old, newStr)
}

func (s *STRINGS) Right(length int) string {
	result := ""
	if length >= len(s.value) {
		result = s.value
	} else {
		result = s.value[len(s.value)-length:]
	}
	return result
}

func (s *STRINGS) Split(sep string) []string {
	return __strings.Split(s.value, sep)
}

func (s *STRINGS) SplitAfter(sep string) []string {
	return __strings.SplitAfter(s.value, sep)
}

func (s *STRINGS) SplitAfterN(sep string, n int) []string {
	return __strings.SplitAfterN(s.value, sep, n)
}

func (s *STRINGS) SplitN(sep string, n int) []string {
	return __strings.SplitN(s.value, sep, n)
}

func (s *STRINGS) String(value ...string) string {
	if len(value) > 0 {
		s.value = value[0]
	}
	return s.value
}

func (s *STRINGS) SubString(start int) string {
	result := ""
	if start < 0 {
		start = 0
	}
	if start >= len(s.value) {
		result = ""
	} else {
		result = s.value[start:]
	}
	return result
}

func (s *STRINGS) Title() string {
	return __strings.ToTitle(s.value)
}

func (s *STRINGS) TitleSpecial(c unicode.SpecialCase) string {
	return __strings.ToTitleSpecial(c, s.value)
}

func (s *STRINGS) Trim(cutset string) string {
	return __strings.Trim(s.value, cutset)
}

func (s *STRINGS) TrimFunc(f func(rune) bool) string {
	return __strings.TrimFunc(s.value, f)
}

func (s *STRINGS) TrimLeft(cutset string) string {
	return __strings.TrimLeft(s.value, cutset)
}

func (s *STRINGS) TrimLeftFunc(f func(rune) bool) string {
	return __strings.TrimLeftFunc(s.value, f)
}

func (s *STRINGS) TrimPrefix(prefix string) string {
	return __strings.TrimPrefix(s.value, prefix)
}

func (s *STRINGS) TrimRight(cutset string) string {
	return __strings.TrimRight(s.value, cutset)
}

func (s *STRINGS) TrimRightFunc(f func(rune) bool) string {
	return __strings.TrimRightFunc(s.value, f)
}

func (s *STRINGS) TrimSpace() string {
	return __strings.TrimSpace(s.value)
}

func (s *STRINGS) TrimSuffix(suffix string) string {
	return __strings.TrimSuffix(s.value, suffix)
}

func (s *STRINGS) Upper() string {
	return __strings.ToUpper(s.value)
}

func (s *STRINGS) UpperSpecial(c unicode.SpecialCase) string {
	return __strings.ToUpperSpecial(c, s.value)
}
