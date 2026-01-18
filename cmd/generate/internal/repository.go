package internal

import (
	"os"
	"path/filepath"
)

func GenerateRepository(name string) error {
	targetDir := filepath.Join("internal/repository", name)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	data := TemplateData{
		ModuleName:  getModuleName(),
		PackageName: extractPackageName(name),
	}

	err = generateFile(filepath.Join(targetDir, "repository.go"), repositoryTemplate, data)
	if err != nil {
		return err
	}

	return nil
}

const repositoryTemplate = `package {{.PackageName}}

import (
	"{{.ModuleName}}/internal/bootstrap"
	"{{.ModuleName}}/internal/db"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type Repository interface {
}

type RepositoryImpl struct {
	sql db.Sql
}

func NewRepository(i do.Injector) (Repository, error) {
	return &RepositoryImpl{
		sql: do.MustInvoke[db.Sql](i),
	}, nil
}
`
