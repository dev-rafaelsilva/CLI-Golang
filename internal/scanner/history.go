package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"github.com/dev-rafaelsilva/secretscan/internal/patterns"
)

type HistoryFinding struct {
	Commit   string
	File     string
	Pattern  string
	Severity patterns.Severity
}

func ScanHistory(root string) ([]HistoryFinding, error) {
	cmd := exec.Command("git", "-C", root, "log", "-p", "--no-color", "--all")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("git não encontrado ou %s não é um repositório git: %w", root, err)
	}

	pats := patterns.Default()
	var findings []HistoryFinding

	currentCommit := ""
	currentFile := ""
	seen := map[string]bool{}

	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)

	for sc.Scan() {
		line := sc.Text()

		switch {
		case strings.HasPrefix(line, "commit "):
			currentCommit = strings.TrimSpace(strings.TrimPrefix(line, "commit "))
			if len(currentCommit) > 7 {
				currentCommit = currentCommit[:7]
			}

		case strings.HasPrefix(line, "+++ b/"):
			currentFile = strings.TrimPrefix(line, "+++ b/")

		case strings.HasPrefix(line, "+++ /dev/null"):
			currentFile = ""

		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			content := strings.TrimPrefix(line, "+")
			for _, p := range pats {
				if p.Regex.MatchString(content) {
					key := currentCommit + currentFile + p.Name
					if seen[key] {
						continue
					}
					seen[key] = true
					findings = append(findings, HistoryFinding{
						Commit:   currentCommit,
						File:     currentFile,
						Pattern:  p.Name,
						Severity: p.Severity,
					})
				}
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return findings, fmt.Errorf("erro ao rodar git log: %w", err)
	}

	return findings, sc.Err()
}