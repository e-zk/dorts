package templates

import (
	"bytes"
	"github.com/e-zk/dorts/templates/functions"
	"path"
	"text/template"
)

func newTemplateFile(tmpPath string) (*template.Template, error) {
	var tmp *template.Template
	tmp, err := template.New(path.Base(tmpPath)).Funcs(functions.GetFuncs()).ParseFiles(tmpPath)
	if err != nil {
		return tmp, err
	}

	return tmp, nil
}

/* process template string */
func process(t *template.Template, vars interface{}) (string, error) {
	var err error = nil
	var tmpBytes bytes.Buffer

	err = t.Execute(&tmpBytes, vars)
	if err != nil {
		return "", err
	}

	return tmpBytes.String(), err
}

/* process template string */
func ProcessString(str string, vars interface{}) (string, error) {
	tmp, err := template.New("tmp").Parse(str)
	if err != nil {
		return "", err
	}

	return process(tmp, vars)
}

/* process template file */
func ProcessFile(tmpPath string, vars interface{}) (string, error) {
	tmp, err := newTemplateFile(tmpPath)
	if err != nil {
		return "", err
	}

	return process(tmp, vars)
}
