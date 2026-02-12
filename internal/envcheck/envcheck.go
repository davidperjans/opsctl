package envcheck

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Result struct {
	Missing []string // keys in example but not in env
	Extra   []string // keys in env but not in example
}

func Check(examplePath, envPath string) (Result, error) {
	exampleKeys, err := readKeys(examplePath)
	if err != nil {
		return Result{}, err
	}

	envKeys, err := readKeys(envPath)
	if err != nil {
		return Result{}, err
	}

	missing := diff(exampleKeys, envKeys)
	extra := diff(envKeys, exampleKeys)

	return Result{
		Missing: missing,
		Extra:   extra,
	}, nil
}

func readKeys(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open %s: %w", path, err)
	}
	defer file.Close()

	keysSet := map[string]struct{}{}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}

		// ignore comments in env file
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// allow: export KEY=value
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		// only consider KEY=VALUE lines
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		if key == "" {
			continue
		}

		keysSet[key] = struct{}{}
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("unable to read %s: %w", path, err)
	}

	// stable order (sort later)
	keys := make([]string, 0, len(keysSet))
	for k := range keysSet {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys, nil
}

func diff(a, b []string) []string {
	bSet := map[string]struct{}{}
	for _, k := range b {
		bSet[k] = struct{}{}
	}

	var out []string
	for _, k := range a {
		if _, ok := bSet[k]; !ok {
			out = append(out, k)
		}
	}
	sortStrings(out)
	return out
}

func sortStrings(s []string) {
	// tiny local sorter to avoid importing sort everywhere
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[i] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}
