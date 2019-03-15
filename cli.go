package util

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
)

func Prompt(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	text = strings.TrimRight(text, "`")
	text = strings.TrimLeft(text, "`")
	if strings.Contains(text, "?") {
		newtext := strings.Split(text, "?")
		text = newtext[0]
	}
	return text
}

// Template is a `flag.Value` for `text.Template` arguments.
// The value of the `Root` field is used as a root template when specified.
type TmplFlag struct {
	Root *template.Template

	Value *template.Template
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TmplFlag) Help() string {
	return "a go template"
}

// Set is flag.Value.Set
func (fv *TmplFlag) Set(v string) error {
	root := fv.Root
	if root == nil {
		root = template.New("")
	}
	t, err := root.New(fmt.Sprintf("%T(%p)", fv, fv)).Parse(v)
	if err == nil {
		fv.Value = t
	}
	return err
}

func (fv *TmplFlag) String() string {
	return fv.Text
}

// Templates is a `flag.Value` for `text.Template` arguments.
// The value of the `Root` field is used as a root template when specified.
type TmplsFlag struct {
	Root *template.Template

	Values []*template.Template
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *TmplsFlag) Help() string {
	return "a go template"
}

// Set is flag.Value.Set
func (fv *TmplsFlag) Set(v string) error {
	root := fv.Root
	if root == nil {
		root = template.New("")
	}
	t, err := root.New(fmt.Sprintf("%T(%p)", fv, fv)).Parse(v)
	if err == nil {
		fv.Texts = append(fv.Texts, v)
		fv.Values = append(fv.Values, t)
	}
	return err
}

func (fv *TmplsFlag) String() string {
	return fmt.Sprint(fv.Texts)
}

// File is a `flag.Value` for file path arguments.
// By default, any errors from os.Stat are returned.
// Alternatively, the value of the `Validate` field is used as a validator when specified.
type FileFlag struct {
	Validate func(os.FileInfo, error) error

	Value string
}

// Set is flag.Value.Set
func (fv *FileFlag) Set(v string) error {
	info, err := os.Stat(v)
	fv.Value = v
	if fv.Validate != nil {
		return fv.Validate(info, err)
	}
	return err
}

func (fv *FileFlag) String() string {
	return fv.Value
}

// FilesFlag is a `flag.Value` for file path arguments.
// By default, any errors from os.Stat are returned.
// Alternatively, the value of the `Validate` field is used as a validator when specified.
type FilesFlag struct {
	Validate func(os.FileInfo, error) error

	Values []string
}

// Set is flag.Value.Set
func (fv *FilesFlag) Set(v string) error {
	info, err := os.Stat(v)
	fv.Values = append(fv.Values, v)
	if fv.Validate != nil {
		return fv.Validate(info, err)
	}
	return err
}

func (fv *FilesFlag) String() string {
	return strings.Join(fv.Values, ",")
}

// EnumFlag is a `flag.Value` for one-of-a-fixed-set string arguments.
// The value of the `Choices` field defines the valid choices.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumFlag struct {
	Choices       []string
	CaseSensitive bool

	Value string
	Text  string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumFlag) Help() string {
	if fv.CaseSensitive {
		return fmt.Sprintf("one of %v (case-sensitive)", fv.Choices)
	}
	return fmt.Sprintf("one of %v", fv.Choices)
}

// Set is flag.Value.Set
func (fv *EnumFlag) Set(v string) error {
	fv.Text = v
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	for _, c := range fv.Choices {
		if equal(c, v) {
			fv.Value = c
			return nil
		}
	}
	return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
}

func (fv *EnumFlag) String() string {
	return fv.Value
}

// EnumFlags is a `flag.Value` for one-of-a-fixed-set string arguments.
// The value of the `Choices` field defines the valid choices.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumFlags struct {
	Choices       []string
	CaseSensitive bool

	Values []string
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumFlags) Help() string {
	if fv.CaseSensitive {
		return fmt.Sprintf("one of %v (case-sensitive)", fv.Choices)
	}
	return fmt.Sprintf("one of %v", fv.Choices)
}

// Set is flag.Value.Set
func (fv *EnumFlags) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	for _, c := range fv.Choices {
		if equal(c, v) {
			fv.Values = append(fv.Values, c)
			fv.Texts = append(fv.Texts, v)
			return nil
		}
	}
	return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
}

func (fv *EnumFlags) String() string {
	return strings.Join(fv.Values, ",")
}

// EnumFlagsCSV is a `flag.Value` for comma-separated EnumFlag arguments.
// The value of the `Choices` field defines the valid choices.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumFlagsCSV struct {
	Choices       []string
	Separator     string
	Accumulate    bool
	CaseSensitive bool

	Values []string
	Texts  []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumFlagsCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	if fv.CaseSensitive {
		return fmt.Sprintf("%q-separated list of values from %v (case-sensitive)", separator, fv.Choices)
	}
	return fmt.Sprintf("%q-separated list of values from %v", separator, fv.Choices)
}

// Set is flag.Value.Set
func (fv *EnumFlagsCSV) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	if !fv.Accumulate {
		fv.Values = fv.Values[:0]
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		var ok bool
		var value string
		for _, c := range fv.Choices {
			if equal(c, part) {
				value = c
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
		}
		fv.Values = append(fv.Values, value)
		fv.Texts = append(fv.Texts, part)
	}
	return nil
}

func (fv *EnumFlagsCSV) String() string {
	return strings.Join(fv.Values, ",")
}

// EnumFlagSet is a `flag.Value` for one-of-a-fixed-set string arguments.
// Only distinct values are returned.
// The value of the `Choices` field defines the valid choices.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumFlagSet struct {
	Choices       []string
	CaseSensitive bool

	Value map[string]bool
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumFlagSet) Help() string {
	if fv.CaseSensitive {
		return fmt.Sprintf("one of %v (case-sensitive)", fv.Choices)
	}
	return fmt.Sprintf("one of %v", fv.Choices)
}

// Values returns a string slice of specified values.
func (fv *EnumFlagSet) Values() (out []string) {
	for v := range fv.Value {
		out = append(out, v)
	}
	sort.Strings(out)
	return
}

// Set is flag.Value.Set
func (fv *EnumFlagSet) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	var ok bool
	for _, c := range fv.Choices {
		if equal(c, v) {
			v = c
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
	}
	if fv.Value == nil {
		fv.Value = make(map[string]bool)
	}
	fv.Value[v] = true
	fv.Texts = append(fv.Texts, v)
	return nil
}

func (fv *EnumFlagSet) String() string {
	return strings.Join(fv.Values(), ",")
}

// EnumFlagSetCSV is a `flag.Value` for comma-separated EnumFlag arguments.
// Only distinct values are returned.
// The value of the `Choices` field defines the valid choices.
// If `Accumulate` is set, the values of all instances of the flag are accumulated.
// The `Separator` field is used instead of the comma when set.
// If `CaseSensitive` is set to `true` (default `false`), the comparison is case-sensitive.
type EnumFlagSetCSV struct {
	Choices       []string
	Separator     string
	Accumulate    bool
	CaseSensitive bool

	Value map[string]bool
	Texts []string
}

// Help returns a string suitable for inclusion in a flag help message.
func (fv *EnumFlagSetCSV) Help() string {
	separator := ","
	if fv.Separator != "" {
		separator = fv.Separator
	}
	if fv.CaseSensitive {
		return fmt.Sprintf("%q-separated list of values from %v (case-sensitive)", separator, fv.Choices)
	}
	return fmt.Sprintf("%q-separated list of values from %v", separator, fv.Choices)
}

// Values returns a string slice of specified values.
func (fv *EnumFlagSetCSV) Values() (out []string) {
	for v := range fv.Value {
		out = append(out, v)
	}
	sort.Strings(out)
	return
}

// Set is flag.Value.Set
func (fv *EnumFlagSetCSV) Set(v string) error {
	equal := strings.EqualFold
	if fv.CaseSensitive {
		equal = func(a, b string) bool { return a == b }
	}
	separator := fv.Separator
	if separator == "" {
		separator = ","
	}
	if !fv.Accumulate || fv.Value == nil {
		fv.Value = make(map[string]bool)
	}
	parts := strings.Split(v, separator)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		var ok bool
		var value string
		for _, c := range fv.Choices {
			if equal(c, part) {
				value = c
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(`"%s" must be one of [%s]`, v, strings.Join(fv.Choices, " "))
		}
		fv.Value[value] = true
		fv.Texts = append(fv.Texts, part)
	}
	return nil
}

func (fv *EnumFlagSetCSV) String() string {
	return strings.Join(fv.Values(), ",")
}
