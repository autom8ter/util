package util

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RenderTmpl(t string, data interface{}) string {
	if strings.Contains(t, "{{") {
		tmpl, err := template.New("").Funcs(sprig.GenericFuncMap()).Parse(t)
		FatalIfErr(err, "failed to create template to render string", t)
		buf := bytes.NewBuffer(nil)
		if err := tmpl.Execute(buf, data); err != nil {
			FatalIfErr(err, "failed to render string at execution", t)
		}
		return buf.String()
	}
	return t
}

func WalkTemplatesFromConfig(dir string, outDir string) {
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Exit(1, errFmt, err, "error walking path")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := fs.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, viper.AllSettings())
		}
		return nil
	}); err != nil {
		Exit(1, errFmt, err, "failed to walk templates")
	}
}

func WalkTemplates(dir, outDir string, data interface{}) {
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Exit(1, errFmt, err, "error walking path")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := fs.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, data)
		}
		return nil
	}); err != nil {
		Exit(1, errFmt, err, "failed to walk templates")
	}
}

type AssetFunc func(string) ([]byte, error)
type AssetDirFunc func(string) ([]string, error)

func loadDirectory(dfn AssetDirFunc, afn AssetFunc, directory string) (*template.Template, error) {
	var tmpl *template.Template

	return assetToTmpl(dfn, afn, tmpl, directory)
}

func loadHtmlDirectory(dfn AssetDirFunc, afn AssetFunc, directory string) (*template.Template, error) {
	var tmpl *template.Template
	return htmlassetToTmpl(dfn, afn, tmpl, directory)
}

func assetToTmpl(dfn AssetDirFunc, afn AssetFunc, tmpl *template.Template, directory string) (*template.Template, error) {
	files, err := dfn(directory)
	if err != nil {
		return tmpl, err
	}

	for _, filePath := range files {
		contents, err := afn(directory + "/" + filePath)
		if err != nil {
			return tmpl, err
		}

		name := filepath.Base(filePath)

		if tmpl == nil {
			tmpl = template.New(name)
		}

		if name != tmpl.Name() {
			tmpl = tmpl.New(name)
		}
		tmpl.Funcs(sprig.GenericFuncMap())

		if _, err = tmpl.Parse(string(contents)); err != nil {
			return tmpl, err
		}
	}

	return tmpl, nil
}

func htmlassetToTmpl(dfn AssetDirFunc, afn AssetFunc, tmpl *template.Template, directory string) (*template.Template, error) {
	files, err := dfn(directory)
	if err != nil {
		return tmpl, err
	}

	for _, filePath := range files {
		contents, err := afn(directory + "/" + filePath)
		if err != nil {
			return tmpl, err
		}

		name := filepath.Base(filePath)

		if tmpl == nil {
			tmpl = template.New(name)
		}

		if name != tmpl.Name() {
			tmpl = tmpl.New(name)
		}
		tmpl.Funcs(sprig.GenericFuncMap())

		if _, err = tmpl.Parse(string(contents)); err != nil {
			return tmpl, err
		}
	}

	return tmpl, nil
}

func MustParseAssets(dfn AssetDirFunc, afn AssetFunc, directory string) *template.Template {
	if tmpl, err := loadDirectory(dfn, afn, directory); err != nil {
		panic(err)
	} else {
		return tmpl
	}
}

func MustParseHtmlAssets(dfn AssetDirFunc, afn AssetFunc, directory string) *template.Template {
	if tmpl, err := loadHtmlDirectory(dfn, afn, directory); err != nil {
		panic(err)
	} else {
		return tmpl
	}
}

// Template reads a go template and writes it to dist given data.
func MustExecAssets(dfn AssetDirFunc, afn AssetFunc, dir string, data interface{}, w io.Writer) error {
	tmpl := MustParseAssets(dfn, afn, dir)
	return tmpl.Execute(w, data)
}

// Template reads a go template and writes it to dist given data.
func MustExecHtmlAssets(dfn AssetDirFunc, afn AssetFunc, dir string, data interface{}, w io.Writer) error {
	tmpl := MustParseHtmlAssets(dfn, afn, dir)
	return tmpl.Execute(w, data)
}

func RenderFilesToResponseWriter(w http.ResponseWriter, relTmplPath string, data interface{}) {
	cwd, _ := os.Getwd()
	t, err := template.ParseFiles(filepath.Join(cwd, relTmplPath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
