package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var makeResponseCmd = &cobra.Command{
	Use:   "make:response [name]",
	Short: "Create a new response",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		makeResponse(name)
	},
}

func init() {
	rootCmd.AddCommand(makeResponseCmd)
}

func makeResponse(name string) {
	// گرفتن مسیر و نام سرویس
	parts := strings.Split(name, "/")
	serviceName := parts[len(parts)-1]
	dir := filepath.Join(append([]string{"app", "http", "responses"}, parts[:len(parts)-1]...)...)

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
	pkgName := "responses"
	if len(parts) > 1 {
		pkgName = parts[len(parts)-2]
	}

	// تولید محتوای فایل
	content := fmt.Sprintf(`package %s

type %s struct {
	Name string 
}
`, pkgName, serviceName)

	// نوشتن فایل
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Fatalf("error writing file: %v\n", err)
		return
	}

	fmt.Println("✅ response created:", filePath)
}
