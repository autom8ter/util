package util

import (
	"bufio"
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	util "github.com/Masterminds/goutils"
	"github.com/Masterminds/sprig"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"unicode"
)

func Base64encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func Base64decode(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func Base32encode(v string) string {
	return base32.StdEncoding.EncodeToString([]byte(v))
}

func Base32decode(v string) string {
	data, err := base32.StdEncoding.DecodeString(v)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func Abbrev(width int, s string) string {
	if width < 4 {
		return s
	}
	r, _ := util.Abbreviate(s, width)
	return r
}

/*
DeleteWhiteSpace deletes all whitespaces from a string as defined by unicode.IsSpace(rune).
It returns the string without whitespaces.
Parameter:
    str - the string to delete whitespace from, may be nil
Returns:
    the string without whitespaces
*/
func DeleteWhiteSpace(str string) string {
	if str == "" {
		return str
	}
	sz := len(str)
	var chs bytes.Buffer
	count := 0
	for i := 0; i < sz; i++ {
		ch := rune(str[i])
		if !unicode.IsSpace(ch) {
			chs.WriteRune(ch)
			count++
		}
	}
	if count == sz {
		return str
	}
	return chs.String()
}

func Abbrevboth(left, right int, s string) string {
	if right < 4 || left > 0 && right < 7 {
		return s
	}
	r, _ := util.AbbreviateFull(s, left, right)
	return r
}
func Initials(s string) string {
	// Wrap this just to eliminate the var args, which templates don't do well.
	return util.Initials(s)
}

func RandAlphaNumeric(count int) string {
	// It is not possible, it appears, to actually generate an error here.
	r, _ := util.CryptoRandomAlphaNumeric(count)
	return r
}

func RandAlpha(count int) string {
	r, _ := util.CryptoRandomAlphabetic(count)
	return r
}

func RandAscii(count int) string {
	r, _ := util.CryptoRandomAscii(count)
	return r
}

func RandNumeric(count int) string {
	r, _ := util.CryptoRandomNumeric(count)
	return r
}

func Untitle(str string) string {
	return util.Uncapitalize(str)
}

func Auote(str ...interface{}) string {
	out := make([]string, 0, len(str))
	for _, s := range str {
		if s != nil {
			out = append(out, fmt.Sprintf("%q", Strval(s)))
		}
	}
	return strings.Join(out, " ")
}

func Squote(str ...interface{}) string {
	out := make([]string, 0, len(str))
	for _, s := range str {
		if s != nil {
			out = append(out, fmt.Sprintf("'%v'", s))
		}
	}
	return strings.Join(out, " ")
}

func Cat(v ...interface{}) string {
	v = RemoveNilElements(v)
	r := strings.TrimSpace(strings.Repeat("%v ", len(v)))
	return fmt.Sprintf(r, v...)
}

func Plural(one, many string, count int) string {
	if count == 1 {
		return one
	}
	return many
}

func Strslice(v interface{}) []string {
	switch v := v.(type) {
	case []string:
		return v
	case []interface{}:
		b := make([]string, 0, len(v))
		for _, s := range v {
			if s != nil {
				b = append(b, Strval(s))
			}
		}
		return b
	default:
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			l := val.Len()
			b := make([]string, 0, l)
			for i := 0; i < l; i++ {
				value := val.Index(i).Interface()
				if value != nil {
					b = append(b, Strval(value))
				}
			}
			return b
		default:
			if v == nil {
				return []string{}
			} else {
				return []string{Strval(v)}
			}
		}
	}
}

func RemoveNilElements(v []interface{}) []interface{} {
	newSlice := make([]interface{}, 0, len(v))
	for _, i := range v {
		if i != nil {
			newSlice = append(newSlice, i)
		}
	}
	return newSlice
}

func Strval(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func Trunc(c int, s string) string {
	if len(s) <= c {
		return s
	}
	return s[0:c]
}

func Join(sep string, v interface{}) string {
	return strings.Join(Strslice(v), sep)
}

func Split(sep, orig string) map[string]string {
	parts := strings.Split(orig, sep)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

func Splitn(sep string, n int, orig string) map[string]string {
	parts := strings.SplitN(orig, sep, n)
	res := make(map[string]string, len(parts))
	for i, v := range parts {
		res["_"+strconv.Itoa(i)] = v
	}
	return res
}

// substring creates a substring of the given string.
//
// If start is < 0, this calls string[:end].
//
// If start is >= 0 and end < 0 or end bigger than s length, this calls string[start:]
//
// Otherwise, this calls string[start, end].
func Substring(start, end int, s string) string {
	if start < 0 {
		return s[:end]
	}
	if end < 0 || end > len(s) {
		return s[start:]
	}
	return s[start:end]
}

// ToPrettyJsonString encodes an item into a pretty (indented) JSON string
func ToPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}

// ToPrettyJson encodes an item into a pretty (indented) JSON
func ToPrettyJson(obj interface{}) []byte {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return output
}

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func ScanAndReplace(r io.Reader, replacements ...string) string {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	var text string
	for scanner.Scan() {
		text = rep.Replace(scanner.Text())
	}
	return text
}

func Render(s string, data interface{}) string {
	if strings.Contains(s, "{{") {
		t, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(s)
		FatalIfErr(err, "failed to create template to render string", s)
		buf := bytes.NewBuffer(nil)
		if err := t.Execute(buf, data); err != nil {
			FatalIfErr(err, "failed to render string at execution", s)
		}
		return buf.String()
	}
	return s
}

func HorizontalLine() {
	fmt.Println("<------------------------------------------------------------------------------------>")
}

func SingleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
