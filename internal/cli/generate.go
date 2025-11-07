package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/platforma-dev/platforma/log"
)

func generateCommand(args []string) {
	subject := args[0]

	switch subject {
	case "domain":
		domainName := args[1]
		log.Info("generating domain", "domain", domainName)

		data := struct {
			PackageName string
			TypeName    string
		}{
			PackageName: domainName,
			TypeName:    title(domainName),
		}

		err := writeFromTemplate("./internal/"+domainName, "model.go", "templates/domain/model.go.tmpl", data)
		if err != nil {
			log.Error("error", "error", err)
			return
		}

		err = writeFromTemplate("./internal/"+domainName, "repository.go", "templates/domain/repository.go.tmpl", data)
		if err != nil {
			log.Error("error", "error", err)
			return
		}

		err = writeFromTemplate("./internal/"+domainName, "service.go", "templates/domain/service.go.tmpl", data)
		if err != nil {
			log.Error("error", "error", err)
			return
		}

		err = writeFromTemplate("./internal/"+domainName, "domain.go", "templates/domain/domain.go.tmpl", data)
		if err != nil {
			log.Error("error", "error", err)
			return
		}
	default:
		log.Error("can't generate subject", "subject", subject)
	}
}

func writeFromTemplate(folder, file, templatePath string, data any) error {
	err := os.MkdirAll(folder, 0750)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", folder, err)
	}

	// Get the directory of the current CLI package
	_, filename, _, _ := runtime.Caller(0)
	cliDir := filepath.Dir(filename)
	fullTemplatePath := filepath.Join(cliDir, templatePath)

	templateContent, err := os.ReadFile(fullTemplatePath) //nolint:gosec // Known path in compile time
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", fullTemplatePath, err)
	}

	tmpl, err := template.New(file).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", fullTemplatePath, err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", fullTemplatePath, err)
	}

	err = os.WriteFile(filepath.Join(folder, file), buf.Bytes(), 0600)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filepath.Join(folder, file), err)
	}

	return nil
}

func title(s string) string {
	return strings.ToUpper(string([]rune(s)[0])) + s[1:]
}
