package internal

import (
	"os"
	"path/filepath"

	"github.com/anonychun/bibit/internal/lib"
)

func GenerateUsecase(name string) error {
	targetDir := filepath.Join("internal/usecase", name)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	data := TemplateData{
		ModuleName:  lib.GetModuleName(),
		PackageName: lib.ExtractPackageName(name),
	}

	err = generateFile(filepath.Join(targetDir, "usecase.go"), usecaseTemplate, data)
	if err != nil {
		return err
	}

	err = generateFile(filepath.Join(targetDir, "http_handler.go"), httpHandlerTemplate, data)
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

const httpHandlerTemplate = `package {{.PackageName}}

import (
	"{{.ModuleName}}/internal/bootstrap"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewHttpHandler)
}

type IHttpHandler interface {
}

type HttpHandler struct {
	usecase IUsecase
}

var _ IHttpHandler = (*HttpHandler)(nil)

func NewHttpHandler(i do.Injector) (*HttpHandler, error) {
	return &HttpHandler{
		usecase: do.MustInvoke[*Usecase](i),
	}, nil
}
`
