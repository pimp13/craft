package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:     "make:module [name]",
	Short:   "Create a new module with boilerplate files",
	Aliases: []string{"g:mod [name]"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateModule(args[0])
	},
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}

type moduleTemplateData struct {
	PackageName string
	StructName  string
	ServiceName string
}

func generateModule(moduleName string) error {
	pathParts := strings.Split(moduleName, "/")
	packageName := strings.ToLower(pathParts[len(pathParts)-1])
	structName := strings.Title(packageName)

	data := moduleTemplateData{
		PackageName: packageName,
		StructName:  structName,
		ServiceName: strings.ToLower(structName),
	}

	basePath := filepath.Join(append([]string{"internal", "module"}, pathParts...)...)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	templates := map[string]string{
		"service.go":    serviceTemplate,
		"controller.go": controllerTemplate,
		"dto.go":        dtoTemplate,
		"wire.go":       wireTemplate,
		"middleware.go": middlewareTemplate,
	}

	for fileName, tmplContent := range templates {
		filePath := filepath.Join(basePath, fileName)

		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("⚠️  file %s already exists, skipping...\n", filePath)
			continue
		}

		tmpl, err := template.New(fileName).Parse(tmplContent)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", fileName, err)
		}

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
		defer file.Close()

		if err := tmpl.Execute(file, data); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", fileName, err)
		}

		fmt.Printf("✅ created %s\n", filePath)
	}

	fmt.Printf("🎉 module %s generated successfully!\n", moduleName)
	return nil
}

const serviceTemplate = `package {{.PackageName}}

import "context"

type {{.StructName}}Service interface {
	Index(ctx context.Context)
}

type {{.StructName}}ServiceImpl struct {}

func New{{.StructName}}Service() {{.StructName}}Service {
	return &{{.StructName}}ServiceImpl{}
}

func (s *{{.StructName}}ServiceImpl) Index(ctx context.Context) {}
`

const controllerTemplate = `package {{.PackageName}}

import "github.com/labstack/echo/v4"

type {{.StructName}}Controller struct {
	{{.ServiceName}}Service {{.StructName}}Service
}

func New{{.StructName}}Controller({{.ServiceName}}Service {{.StructName}}Service) *{{.StructName}}Controller {
	return &{{.StructName}}Controller{
		{{.ServiceName}}Service,
	}
}

func (ctrl *{{.StructName}}Controller) Routes(r *echo.Group) {
	r.GET("/{{.ServiceName}}", ctrl.index)
}

func (ctrl *{{.StructName}}Controller) index(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "Hello World",
		"ok":      true,
	})
}
`

const dtoTemplate = `package {{.PackageName}}

type {{.StructName}}Dto struct {
	Name string ` + "`json:\"name\"`" + `
}
`

const wireTemplate = `//go:build wireinject
// +build wireinject

package {{.PackageName}}

import "github.com/google/wire"

var Provider = wire.NewSet(
	New{{.StructName}}Service,
	New{{.StructName}}Controller,
	New{{.StructName}}Middleware,
)
`

const middlewareTemplate = `package {{.PackageName}}

type {{.StructName}}Middleware struct {}

func New{{.StructName}}Middleware() *{{.StructName}}Middleware {
	return &{{.StructName}}Middleware{}
}
`
