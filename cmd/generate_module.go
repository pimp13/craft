package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
	Use:     "generate:module [name]",
	Short:   "Generate a new module with boilerplate files",
	Aliases: []string{"g:mod [name]"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := args[0]
		return generateModule(moduleName)
	},
}

func generateModule(moduleName string) error {
	packageName := strings.ToLower(moduleName)
	structName := strings.Title(moduleName)

	basePath := filepath.Join("internal", "module", packageName)
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create module directory: %w", err)
	}

	// محتواها
	contents := map[string]string{
		"service.go": fmt.Sprintf(`package %s

import "context"

type %sService interface {
  Index(ctx context.Context)
}

type %sServiceImpl struct {}

func New%sService() %sService {
  return &%sServiceImpl{}
}

func (s *%sServiceImpl) Index(ctx context.Context) {}
`, packageName, structName, structName, structName, structName, structName, structName),

		"controller.go": fmt.Sprintf(`package %s

import "github.com/labstack/echo/v4"

type %sController struct {
  %sService %sService
}

func New%sController(%sService %sService) *%sController {
  return &%sController{
    %sService,
  }
}

func (s *%sController) Index(c echo.Context) {}
`, packageName, structName, strings.ToLower(structName), structName,
			structName, strings.ToLower(structName), structName,
			structName, structName,
			strings.ToLower(structName),
			structName),

		"dto.go": fmt.Sprintf(`package %s

type %sDto struct {
  Name string `+"`json:\"name\"`"+`
}
`, packageName, structName),

		"wire.go": fmt.Sprintf(`//go:build wireinject
// +build wireinject

package %s

import "github.com/google/wire"

var Provider = wire.NewSet(
  New%sService,
  New%sController,
  New%sMiddleware,
)
`, packageName, structName, structName, structName),

		"middleware.go": fmt.Sprintf(`package %s

type %sMiddleware struct {}

func New%sMiddleware() *%sMiddleware {
  return &%sMiddleware{}
}
`, packageName, structName, structName, structName, structName),
	}

	// ساخت فایل‌ها
	for name, content := range contents {
		filePath := filepath.Join(basePath, name)

		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("⚠️  file %s already exists, skipping...\n", filePath)
			continue
		}

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
		fmt.Printf("✅ created %s\n", filePath)
	}

	fmt.Printf("🎉 module %s generated successfully!\n", moduleName)
	return nil
}

func init() {
	rootCmd.AddCommand(moduleCmd)
}
