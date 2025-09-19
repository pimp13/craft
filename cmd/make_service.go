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
		makeService(name)
	},
}

func init() {
	rootCmd.AddCommand(makeServiceCmd)
}

func makeService(name string) {
	// گرفتن مسیر و نام سرویس
	parts := strings.Split(name, "/")
	serviceName := parts[len(parts)-1]
	dir := filepath.Join(append([]string{"app", "services"}, parts[:len(parts)-1]...)...)

	// ساخت پوشه اگر وجود نداشت
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
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
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Fatalf("error writing file: %v\n", err)
		return
	}

	fmt.Println("✅ service created:", filePath)
}
