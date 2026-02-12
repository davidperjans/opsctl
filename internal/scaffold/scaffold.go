package scaffold

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	ServiceName string
	ModulePath  string
	TargetDir   string
	Force       bool
}

func GenerateFromFS(templateFS fs.FS, templateRoot string, opts Options) error {
	if opts.ServiceName == "" {
		return errors.New("service name is required")
	}
	if opts.ModulePath == "" {
		return errors.New("module path is required")
	}
	if opts.TargetDir == "" {
		return errors.New("target dir is required")
	}

	// ensure target dir
	if _, err := os.Stat(opts.TargetDir); err == nil {
		if !opts.Force {
			return fmt.Errorf("target directory already exists: %s (use --force to overwrite)", opts.TargetDir)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(opts.TargetDir, 0755); err != nil {
			return err
		}
	} else {
		return err
	}

	return fs.WalkDir(templateFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// read template file
		b, err := fs.ReadFile(templateFS, path)
		if err != nil {
			return err
		}

		rel := strings.TrimPrefix(path, templateRoot)
		rel = strings.TrimPrefix(rel, "/")

		// drop .tmpl suffix
		outRel := strings.TrimSuffix(rel, ".tmpl")
		outPath := filepath.Join(opts.TargetDir, filepath.FromSlash(outRel))

		// ensure directory exists
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		// token replace (simple, no templating engine)
		content := string(b)
		content = strings.ReplaceAll(content, "{{SERVICE_NAME}}", opts.ServiceName)
		content = strings.ReplaceAll(content, "{{MODULE_PATH}}", opts.ModulePath)

		// write file
		if !opts.Force {
			if _, err := os.Stat(outPath); err == nil {
				return fmt.Errorf("file already exists: %s (use --force)", outPath)
			}
		}

		return os.WriteFile(outPath, []byte(content), 0644)
	})
}
