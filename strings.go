// --------------------------------------------------------------------------------
// File:        strings.go
// Author:      TRAE AI
// Created:     2025/12/20 12:31:58
// Description: STRING is a wrapper for string operations, providing a
//              convenient way to perform string manipulations in Go.
// --------------------------------------------------------------------------------

package boost

import (
	"strings"
	"unicode"
)

// STRING provides utility methods for string operations.
type STRING struct {
	s string
}

// Compare compares the string with another string lexicographically.
// t: String to compare with
// Returns: -1 if s < t, 0 if s == t, 1 if s > t
// Usage:
// result := STRING{s: "a"}.Compare("b")
// returns -1
func (s *STRING) Compare(t string) int {
	return strings.Compare(s.s, t)
}

// Contains checks if the string contains the specified substring.
// substr: Substring to check for
// Returns: true if substring is found, false otherwise
// Usage:
// hasSubstring := STRING{s: "hello"}.Contains("ell")
// returns true
func (s *STRING) Contains(substr string) bool {
	return strings.Contains(s.s, substr)
}

// ContainsAny checks if the string contains any of the specified characters.
// chars: Characters to check for
// Returns: true if any character is found, false otherwise
// Usage:
// hasAny := STRING{s: "hello"}.ContainsAny("union")
// returns true
func (s *STRING) ContainsAny(chars string) bool {
	return strings.ContainsAny(s.s, chars)
}

// ContainsRune checks if the string contains the specified rune.
// r: Rune to check for
// Returns: true if rune is found, false otherwise
// Usage:
// hasRune := STRING{s: "hello"}.ContainsRune('e')
// returns true
func (s *STRING) ContainsRune(r rune) bool {
	return strings.ContainsRune(s.s, r)
}

// Count counts the number of non-overlapping occurrences of the substring.
// substr: Substring to count
// Returns: Number of occurrences
// Usage:
// count := STRING{s: "hello"}.Count("l")
// returns 2
func (s *STRING) Count(substr string) int {
	return strings.Count(s.s, substr)
}

// Cut splits the string into two parts at the first occurrence of sep.
// sep: Separator string
// Returns: Before separator, after separator, and whether separator was found
// Usage:
// before, after, found := STRING{s: "a/b/c"}.Cut("/")
// "a", "b/c", true
func (s *STRING) Cut(sep string) (before, after STRING, found bool) {
	b, a, f := strings.Cut(s.s, sep)
	return STRING{s: b}, STRING{s: a}, f
}

// CutPrefix removes the specified prefix if present.
// prefix: Prefix to remove
// Returns: String with prefix removed, and whether prefix was found
// Usage:
// result, found := STRING{s: "hello"}.CutPrefix("he")
// "llo", true
func (s *STRING) CutPrefix(prefix string) (STRING, bool) {
	result, found := strings.CutPrefix(s.s, prefix)
	return STRING{s: result}, found
}

// CutSuffix removes the specified suffix if present.
// suffix: Suffix to remove
// Returns: String with suffix removed, and whether suffix was found
// Usage:
// result, found := STRING{s: "hello"}.CutSuffix("lo")
// "hel", true
func (s *STRING) CutSuffix(suffix string) (STRING, bool) {
	result, found := strings.CutSuffix(s.s, suffix)
	return STRING{s: result}, found
}

// EqualFold checks if the string is equal to another string, ignoring case.
// t: String to compare with
// Returns: true if strings are equal ignoring case, false otherwise
// Usage:
// equal := STRING{s: "HELLO"}.EqualFold("hello")
// returns true
func (s *STRING) EqualFold(t string) bool {
	return strings.EqualFold(s.s, t)
}

// Fields splits the string into whitespace-separated fields.
// Returns: Slice of fields
// Usage:
// fields := STRING{s: "hello world  test"}.Fields()
// ["hello", "world", "test"]
func (s *STRING) Fields() []string {
	return strings.Fields(s.s)
}

// FieldsFunc splits the string into fields using the specified function to determine separators.
// f: Function that returns true for separator runes
// Returns: Slice of fields
// Usage:
// isSeparator := func(r rune) bool { return r == ',' }
// fields := STRING{s: "a,b,c"}.FieldsFunc(isSeparator)
// ["a", "b", "c"]
func (s *STRING) FieldsFunc(f func(rune) bool) []string {
	return strings.FieldsFunc(s.s, f)
}

// HasPrefix checks if the string starts with the specified prefix.
// prefix: Prefix to check for
// Returns: true if string starts with prefix, false otherwise
// Usage:
// hasPrefix := STRING{s: "hello"}.HasPrefix("he")
// returns true
func (s *STRING) HasPrefix(prefix string) bool {
	return strings.HasPrefix(s.s, prefix)
}

// HasSuffix checks if the string ends with the specified suffix.
// suffix: Suffix to check for
// Returns: true if string ends with suffix, false otherwise
// Usage:
// hasSuffix := STRING{s: "hello"}.HasSuffix("lo")
// returns true
func (s *STRING) HasSuffix(suffix string) bool {
	return strings.HasSuffix(s.s, suffix)
}

// Index returns the index of the first occurrence of the substring.
// substr: Substring to find
// Returns: Index of first occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.Index("ell")
// returns 1
func (s *STRING) Index(substr string) int {
	return strings.Index(s.s, substr)
}

// IndexAny returns the index of the first occurrence of the specified characters.
// chars: Characters to find
// Returns: Index of first occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.IndexAny("union")
// returns 1
func (s *STRING) IndexAny(chars string) int {
	return strings.IndexAny(s.s, chars)
}

// IndexByte returns the index of the first occurrence of the specified byte.
// c: Byte to find
// Returns: Index of first occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.IndexByte('e')
// returns 1
func (s *STRING) IndexByte(c byte) int {
	return strings.IndexByte(s.s, c)
}

// IndexFunc returns the index of the first rune that satisfies the function.
// f: Function that returns true for the desired rune
// Returns: Index of first matching rune, or -1 if not found
// Usage:
// isVowel := func(r rune) bool { return strings.ContainsRune("union", r) }
// index := STRING{s: "hello"}.IndexFunc(isVowel)
// returns 1
func (s *STRING) IndexFunc(f func(rune) bool) int {
	return strings.IndexFunc(s.s, f)
}

// IndexRune returns the index of the first occurrence of the specified rune.
// r: Rune to find
// Returns: Index of first occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.IndexRune('e')
// returns 1
func (s *STRING) IndexRune(r rune) int {
	return strings.IndexRune(s.s, r)
}

// Join joins the elements of a slice into a single string using the string as separator.
// elems: Slice of strings to join
// Returns: Joined string
// Usage:
// joined := STRING{s: ", "}.Join([]string{"a", "b", "c"})
// returns "a, b, c"
func (s *STRING) Join(elems []string) STRING {
	return STRING{s: strings.Join(elems, s.s)}
}

// Left returns the leftmost n characters of the string.
// length: Number of characters to return from the left
// Returns: Leftmost n characters
// Usage:
// left := STRING{s: "hello"}.Left(3)
// returns "hel"
func (s *STRING) Left(length int) STRING {
	str := s.s
	var result STRING
	if length >= len(str) {
		result = STRING{s: str}
	} else {
		result = STRING{s: str[:length]}
	}
	return result
}

// Mid returns a substring starting at the specified index with the given length.
// start: Starting index (0-based)
// length: Number of characters to return
// Returns: Substring from start index with specified length
// Usage:
// mid := STRING{s: "hello"}.Mid(1, 3)
// returns "ell"
func (s *STRING) Mid(start, length int) STRING {
	str := s.s
	var result STRING
	if start < 0 {
		start = 0
	}
	if start >= len(str) {
		result = STRING{s: ""}
	} else {
		end := start + length
		if end > len(str) {
			end = len(str)
		}
		result = STRING{s: str[start:end]}
	}
	return result
}

// Right returns the rightmost n characters of the string.
// length: Number of characters to return from the right
// Returns: Rightmost n characters
// Usage:
// right := STRING{s: "hello"}.Right(3)
// returns "llo"
func (s *STRING) Right(length int) STRING {
	str := s.s
	var result STRING
	if length >= len(str) {
		result = STRING{s: str}
	} else {
		result = STRING{s: str[len(str)-length:]}
	}
	return result
}

// SubStr returns a substring starting at the specified index.
// start: Starting index (0-based)
// Returns: Substring from start index to end of string
// Usage:
// substr := STRING{s: "hello"}.SubStr(2)
// returns "llo"
func (s *STRING) SubStr(start int) STRING {
	str := s.s
	var result STRING
	if start < 0 {
		start = 0
	}
	if start >= len(str) {
		result = STRING{s: ""}
	} else {
		result = STRING{s: str[start:]}
	}
	return result
}

// LastIndex returns the index of the last occurrence of the substring.
// substr: Substring to find
// Returns: Index of last occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.LastIndex("l")
// returns 3
func (s *STRING) LastIndex(substr string) int {
	return strings.LastIndex(s.s, substr)
}

// LastIndexAny returns the index of the last occurrence of the specified characters.
// chars: Characters to find
// Returns: Index of last occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.LastIndexAny("union")
// returns 4
func (s *STRING) LastIndexAny(chars string) int {
	return strings.LastIndexAny(s.s, chars)
}

// LastIndexByte returns the index of the last occurrence of the specified byte.
// c: Byte to find
// Returns: Index of last occurrence, or -1 if not found
// Usage:
// index := STRING{s: "hello"}.LastIndexByte('l')
// returns 3
func (s *STRING) LastIndexByte(c byte) int {
	return strings.LastIndexByte(s.s, c)
}

// LastIndexFunc returns the index of the last rune that satisfies the function.
// f: Function that returns true for the desired rune
// Returns: Index of last matching rune, or -1 if not found
// Usage:
// isVowel := func(r rune) bool { return strings.ContainsRune("union", r) }
// index := STRING{s: "hello"}.LastIndexFunc(isVowel)
// returns 4
func (s *STRING) LastIndexFunc(f func(rune) bool) int {
	return strings.LastIndexFunc(s.s, f)
}

// Map applies the specified mapping function to each rune in the string.
// mapping: Function to apply to each rune
// Returns: New string with mapping applied
// Usage:
// mapped := STRING{s: "hello"}.Map(func(r rune) rune { return r + 1 })
// returns "ifmmp"
func (s *STRING) Map(mapping func(rune) rune) STRING {
	return STRING{s: strings.Map(mapping, s.s)}
}

// Repeat returns a new string consisting of the original string repeated count times.
// count: Number of times to repeat
// Returns: Repeated string
// Usage:
// repeated := STRING{s: "ab"}.Repeat(3)
// returns "ababab"
func (s *STRING) Repeat(count int) STRING {
	return STRING{s: strings.Repeat(s.s, count)}
}

// Replace replaces the first n occurrences of old substring with new substring.
// old: Substring to replace
// new: Substring to replace with
// n: Maximum number of replacements
// Returns: String with replacements applied
// Usage:
// replaced := STRING{s: "hello hello"}.Replace("hello", "hi", 1)
// returns "hi hello"
func (s *STRING) Replace(old string, new string, n int) STRING {
	return STRING{s: strings.Replace(s.s, old, new, n)}
}

// ReplaceAll replaces all occurrences of old substring with new substring.
// old: Substring to replace
// new: Substring to replace with
// Returns: String with all replacements applied
// Usage:
// replaced := STRING{s: "hello"}.ReplaceAll("hello", "hi")
// returns "hi"
func (s *STRING) ReplaceAll(old string, new string) STRING {
	return STRING{s: strings.ReplaceAll(s.s, old, new)}
}

// Split splits the string into a slice of strings separated by the specified separator.
// sep: Separator string
// Returns: Slice of strings
// Usage:
// split := STRING{s: "a/b/c"}.Split("/")
// returns ["a", "b", "c"]
func (s *STRING) Split(sep string) []string {
	return strings.Split(s.s, sep)
}

// SplitAfter splits the string after each occurrence of the separator.
// sep: Separator string
// Returns: Slice of strings
// Usage:
// split := STRING{s: "a/b/c"}.SplitAfter("/")
// returns ["a/", "b/", "c"]
func (s *STRING) SplitAfter(sep string) []string {
	return strings.SplitAfter(s.s, sep)
}

// SplitAfterN splits the string after each occurrence of the separator, up to n parts.
// sep: Separator string
// n: Maximum number of parts
// Returns: Slice of strings with at most n parts
// Usage:
// split := STRING{s: "a/b/c"}.SplitAfterN("/", 2)
// returns ["a/", "b/c"]
func (s *STRING) SplitAfterN(sep string, n int) []string {
	return strings.SplitAfterN(s.s, sep, n)
}

// SplitN splits the string into a slice of strings separated by the specified separator, up to n parts.
// sep: Separator string
// n: Maximum number of parts
// Returns: Slice of strings with at most n parts
// Usage:
// split := STRING{s: "a/b/c"}.SplitN("/", 2)
// returns ["a", "b/c"]
func (s *STRING) SplitN(sep string, n int) []string {
	return strings.SplitN(s.s, sep, n)
}

// Title returns a copy of the string with the first character of each word capitalized.
// Returns: Title-cased string
// Usage:
// title := STRING{s: "hello world"}.Title()
// returns "Hello World"
func (s *STRING) Title() STRING {
	return STRING{s: strings.ToTitle(s.s)}
}

// ToLower returns the string converted to lowercase.
// Returns: Lowercase string
// Usage:
// lower := STRING{s: "HELLO"}.ToLower()
// returns "hello"
func (s *STRING) ToLower() STRING {
	return STRING{s: strings.ToLower(s.s)}
}

// ToLowerSpecial returns the string converted to lowercase using the specified special case.
// c: Unicode special case for lowercase conversion
// Returns: Lowercase string
// Usage:
// lower := STRING{s: "HELLO"}.ToLowerSpecial(unicode.TurkishCase)
// returns "hello" with Turkish specific lowercase rules
func (s *STRING) ToLowerSpecial(c unicode.SpecialCase) STRING {
	return STRING{s: strings.ToLowerSpecial(c, s.s)}
}

// ToTitleSpecial returns the string converted to title case using the specified special case.
// c: Unicode special case for title case conversion
// Returns: Title-cased string
// Usage:
// title := STRING{s: "hello world"}.ToTitleSpecial(unicode.TurkishCase)
// returns "HELLO WORLD" with Turkish specific title case rules
func (s *STRING) ToTitleSpecial(c unicode.SpecialCase) STRING {
	return STRING{s: strings.ToTitleSpecial(c, s.s)}
}

// ToUpper returns the string converted to uppercase.
// Returns: Uppercase string
// Usage:
// upper := STRING{s: "hello"}.ToUpper()
// returns "HELLO"
func (s *STRING) ToUpper() STRING {
	return STRING{s: strings.ToUpper(s.s)}
}

// ToUpperSpecial returns the string converted to uppercase using the specified special case.
// c: Unicode special case for uppercase conversion
// Returns: Uppercase string
// Usage:
// upper := STRING{s: "hello"}.ToUpperSpecial(unicode.TurkishCase)
// returns "HELLO" with Turkish specific uppercase rules
func (s *STRING) ToUpperSpecial(c unicode.SpecialCase) STRING {
	return STRING{s: strings.ToUpperSpecial(c, s.s)}
}

// Trim removes the leading and trailing characters specified in cutset.
// cutset: Set of characters to remove
// Returns: String with leading and trailing characters removed
// Usage:
// trimmed := STRING{s: "  hello  "}.Trim(" ")
// returns "hello"
func (s *STRING) Trim(cutset string) STRING {
	return STRING{s: strings.Trim(s.s, cutset)}
}

// TrimFunc removes the leading and trailing runes that satisfy the function.
// f: Function that returns true for runes to remove
// Returns: String with leading and trailing runes removed
// Usage:
// trimSpace := func(r rune) bool { return unicode.IsSpace(r) }
// trimmed := STRING{s: "  hello  "}.TrimFunc(trimSpace)
// returns "hello"
func (s *STRING) TrimFunc(f func(rune) bool) STRING {
	return STRING{s: strings.TrimFunc(s.s, f)}
}

// TrimLeft removes the leading characters specified in cutset.
// cutset: Set of characters to remove from left
// Returns: String with leading characters removed
// Usage:
// trimmed := STRING{s: "  hello  "}.TrimLeft(" ")
// returns "hello  "
func (s *STRING) TrimLeft(cutset string) STRING {
	return STRING{s: strings.TrimLeft(s.s, cutset)}
}

// TrimLeftFunc removes the leading runes that satisfy the function.
// f: Function that returns true for runes to remove from left
// Returns: String with leading runes removed
// Usage:
// trimSpace := func(r rune) bool { return unicode.IsSpace(r) }
// trimmed := STRING{s: "  hello  "}.TrimLeftFunc(trimSpace)
// returns "hello  "
func (s *STRING) TrimLeftFunc(f func(rune) bool) STRING {
	return STRING{s: strings.TrimLeftFunc(s.s, f)}
}

// TrimPrefix removes the specified prefix if present.
// prefix: Prefix to remove
// Returns: String with prefix removed
// Usage:
// trimmed := STRING{s: "hello"}.TrimPrefix("he")
// returns "llo"
func (s *STRING) TrimPrefix(prefix string) STRING {
	return STRING{s: strings.TrimPrefix(s.s, prefix)}
}

// TrimRight removes the trailing characters specified in cutset.
// cutset: Set of characters to remove from right
// Returns: String with trailing characters removed
// Usage:
// trimmed := STRING{s: "  hello  "}.TrimRight(" ")
// returns "  hello"
func (s *STRING) TrimRight(cutset string) STRING {
	return STRING{s: strings.TrimRight(s.s, cutset)}
}

// TrimRightFunc removes the trailing runes that satisfy the function.
// f: Function that returns true for runes to remove from right
// Returns: String with trailing runes removed
// Usage:
// trimSpace := func(r rune) bool { return unicode.IsSpace(r) }
// trimmed := STRING{s: "  hello  "}.TrimRightFunc(trimSpace)
// returns "  hello"
func (s *STRING) TrimRightFunc(f func(rune) bool) STRING {
	return STRING{s: strings.TrimRightFunc(s.s, f)}
}

// TrimSpace removes the leading and trailing whitespace from the string.
// Returns: String with whitespace removed
// Usage:
// trimmed := STRING{s: "  hello  "}.TrimSpace()
// returns "hello"
func (s *STRING) TrimSpace() STRING {
	return STRING{s: strings.TrimSpace(s.s)}
}

// TrimSuffix removes the specified suffix if present.
// suffix: Suffix to remove
// Returns: String with suffix removed
// Usage:
// trimmed := STRING{s: "hello"}.TrimSuffix("lo")
// returns "hel"
func (s *STRING) TrimSuffix(suffix string) STRING {
	return STRING{s: strings.TrimSuffix(s.s, suffix)}
}

// String returns the underlying string value or sets a new value.
// If a value is provided, it sets the string to that value and returns it.
// If no value is provided, it returns the current string value.
// Returns: The current string value (after setting if a value was provided)
// Usage:
// str := STRING{s: "hello"}.String()
// returns "hello"
// newStr := str.String("world")
// returns "world" and sets the string to "world"
func (s *STRING) String(value ...string) string {
	if len(value) > 0 {
		s.s = value[0]
	}
	return s.s
}
