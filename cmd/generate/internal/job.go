package internal

import (
	"os"
	"path/filepath"

	"github.com/anonychun/bibit/internal/util"
)

func GenerateJob(name string) error {
	targetDir := filepath.Join("internal/job", name)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	data := TemplateData{
		ModuleName:  util.GetModuleName(),
		PackageName: util.ExtractPackageName(name),
		CmdArg:      name,
	}

	return generateFile(filepath.Join(targetDir, "job.go"), jobTemplate, data)
}

const jobTemplate = `package {{.PackageName}}

import (
	"context"

	"{{.ModuleName}}/internal/bootstrap"
	"github.com/riverqueue/river"
	"github.com/samber/do/v2"
)

func init() {
	do.Provide(bootstrap.Injector, NewJob)
}

type Args struct {
}

func (Args) Kind() string {
	return "{{.CmdArg}}"
}

type Job struct {
	river.WorkerDefaults[Args]
}

func NewJob(i do.Injector) (*Job, error) {
	return &Job{}, nil
}

func (j *Job) Work(ctx context.Context, job *river.Job[Args]) error {
	return nil
}
`
