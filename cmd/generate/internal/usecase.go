package internal

import (
	"os"
	"path/filepath"
)

func GenerateUsecase(name string) error {
	targetDir := filepath.Join("internal/usecase", name)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	data := TemplateData{
		ModuleName:  getModuleName(),
		PackageName: extractPackageName(name),
	}

	err = generateFile(filepath.Join(targetDir, "usecase.go"), usecaseTemplate, data)
	if err != nil {
		return err
	}

	err = generateFile(filepath.Join(targetDir, "handler.go"), handlerTemplate, data)
	if err != nil {
		return err
	}

	err = generateFile(filepath.Join(targetDir, "dto.go"), emptyTemplate, data)
	if err != nil {
		return err
	}

	return nil
}

const usecaseTemplate = `package {{.PackageName}}

import (
	"{{.ModuleName}}/internal/bootstrap"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewUsecase)
}

type IUsecase interface {
}

type Usecase struct {
}

var _ IUsecase = (*Usecase)(nil)

func NewUsecase(i do.Injector) (*Usecase, error) {
	return &Usecase{}, nil
}
`

const handlerTemplate = `package {{.PackageName}}

import (
	"{{.ModuleName}}/internal/bootstrap"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHandler)
}

type IHandler interface {
}

type Handler struct {
	usecase IUsecase
}

var _ IHandler = (*Handler)(nil)

func NewHandler(i do.Injector) (*Handler, error) {
	return &Handler{
		usecase: do.MustInvoke[*Usecase](i),
	}, nil
}
`
