package internal

import (
	"html/template"
	"os"
)

type TemplateData struct {
	CmdArg      string
	ModuleName  string
	PackageName string
}

const emptyTemplate = `package {{.PackageName}}
`

func generateFile(filePath, tmplContent string, data TemplateData) error {
	tmpl, err := template.New("file").Parse(tmplContent)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
