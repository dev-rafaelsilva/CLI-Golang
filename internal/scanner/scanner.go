package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/dev-rafaelsilva/secretscan/internal/patterns"
)

type Finding struct {
	File     string
	Line     int
	Pattern  string
	Severity patterns.Severity
}

var ignoredDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true,
	"dist": true, "build": true, ".venv": true,
}

var ignoredExt = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".ico": true, ".woff": true, ".woff2": true, ".ttf": true,
	".zip": true, ".tar": true, ".gz": true, ".pdf": true,
	".exe": true, ".bin": true, ".so": true,
}

func Scan(root string) ([]Finding, error) {
	var findings []Finding
	pats := patterns.Default()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if ignoredDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if ignoredExt[strings.ToLower(filepath.Ext(path))] {
			return nil
		}

		fileFindings, ferr := scanFile(path, pats)
		if ferr != nil {
			return nil
		}
		findings = append(findings, fileFindings...)
		return nil
	})

	return findings, err
}

func scanFile(path string, pats []patterns.Pattern) ([]Finding, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var findings []Finding
	sc := bufio.NewScanner(f)
	lineNum := 0

	for sc.Scan() {
		lineNum++
		line := sc.Text()
		for _, p := range pats {
			if p.Regex.MatchString(line) {
				findings = append(findings, Finding{
					File: path, Line: lineNum,
					Pattern: p.Name, Severity: p.Severity,
				})
			}
		}
	}
	return findings, sc.Err()
}

func GitignoreMissing(root string) []string {
	sensitivePatterns := []string{".env", "*.pem", "*.key", "*.p12"}

	content := ""
	if data, err := os.ReadFile(filepath.Join(root, ".gitignore")); err == nil {
		content = string(data)
	}

	var missing []string
	for _, sp := range sensitivePatterns {
		if !strings.Contains(content, sp) {
			missing = append(missing, sp)
		}
	}
	return missing
}