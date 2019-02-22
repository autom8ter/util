package util

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"github.com/Masterminds/sprig"
	"crypto/rand"
)

// toPrettyJson encodes an item into a pretty (indented) JSON string
func ToPrettyJsonString(obj interface{}) string {
	output, _ := json.MarshalIndent(obj, "", "  ")
	return fmt.Sprintf("%s", output)
}
// toPrettyJson encodes an item into a pretty (indented) JSON string
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

func Prompt(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("string | " + question)
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

func ScanAndReplace(r io.Reader, replacements ...string) string {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	var text string
	for scanner.Scan() {
		text = rep.Replace(scanner.Text())
	}
	return text
}

func  Render(s string, data interface{}) string {
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

func FatalIfErr(e error, msg string, arg interface{}) {
	if e != nil {
		log.Fatalf("Error: %v Msg: %v Arg: %v", e, msg, arg)
	}
}

func PrintIfErr(e error, msg string, arg interface{}) {
	if e != nil {
		log.Printf("Error: %v Msg: %v Arg: %v", e, msg, arg)
	}
}

func RandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func UserPassword(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Uuid() string {
	return fmt.Sprintf("%s", uuid.New())
}

func Indent(spaces int, v string) string {
	pad := strings.Repeat(" ", spaces)
	return pad + strings.Replace(v, "\n", "\n"+pad, -1)
}

func Replace(old, new, src string) string {
	return strings.Replace(src, old, new, -1)
}