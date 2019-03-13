package netutil

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/autom8ter/util"
	"html/template"
	"strings"
)

type Template string

func (t Template) String() string {
	return string(t)
}

func (t Template) Render(data interface{}) string {
	if strings.Contains(t.String(), "{{") {
		tmpl, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(t.String())
		util.FatalIfErr(err, "failed to create template to render string", t.String())
		buf := bytes.NewBuffer(nil)
		if err := tmpl.Execute(buf, data); err != nil {
			util.FatalIfErr(err, "failed to render string at execution", t)
		}
		return buf.String()
	}
	return t.String()
}

func AuthUserUserTemplate() Template {
	return Template(auth0UserHTML)
}
