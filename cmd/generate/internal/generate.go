package internal

import (
	"html/template"
	"os"
	"runtime/debug"
	"strings"
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

func extractPackageName(name string) string {
	parts := strings.Split(strings.TrimSpace(name), "/")
	return parts[len(parts)-1]
}

func getModuleName() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	if len(bi.Deps) > 0 && bi.Deps[0].Path != "" {
		return bi.Deps[0].Path
	}

	return bi.Path
}
