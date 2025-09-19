package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var makeServiceCmd = &cobra.Command{
	Use:   "make:service [name]",
	Short: "Create a new service",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		createService(name)
	},
}

func init() {
	rootCmd.AddCommand(makeServiceCmd)
}

func createService(name string) {
	// گرفتن مسیر و نام سرویس
	parts := strings.Split(name, "/")
	serviceName := parts[len(parts)-1]
	dir := filepath.Join(append([]string{"app", "services"}, parts[:len(parts)-1]...)...)

	// ساخت پوشه اگر وجود نداشت
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("error creating directory:", err)
		return
	}

	// تبدیل نام سرویس به snake_case
	fileName := toSnakeCase(serviceName) + ".go"
	filePath := filepath.Join(dir, fileName)

	// بررسی وجود فایل
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("file already exists:", filePath)
		return
	}

	// تعیین package name
	pkgName := "services"
	if len(parts) > 1 {
		pkgName = parts[len(parts)-2]
	}

	// تولید محتوای فایل
	content := fmt.Sprintf(`package %s

import (
	"context"

	"github.com/goravel/framework/contracts/database/orm"
)

type %s interface {
	FindAll(ctx context.Context) error
}

type %s struct {
	orm orm.Orm
}

func New%s(orm orm.Orm) %s {
	return &%s{
		orm,
	}
}

func (s *%s) FindAll(ctx context.Context) error {
	return nil
}
`, pkgName, serviceName, toLowerFirst(serviceName), serviceName, serviceName, toLowerFirst(serviceName), toLowerFirst(serviceName))

	// نوشتن فایل
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Fatalf("error writing file: %v\n", err)
		return
	}

	log.Println("✅ service created:", filePath)
}

// utils
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func toLowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
