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
	"context"

	"{{.ModuleName}}/internal/bootstrap"
	dbSql "{{.ModuleName}}/internal/db/sql"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewRepository)
}

type IRepository interface {
}

type Repository struct {
	sqlDB dbSql.IDB
}

var _ IRepository = (*Repository)(nil)

func NewRepository(i do.Injector) (*Repository, error) {
	return &Repository{
		sqlDB: do.MustInvoke[*dbSql.PostgresDB](i),
	}, nil
}
`
