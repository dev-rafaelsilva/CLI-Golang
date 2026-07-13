// Package report formata os achados do scanner para saída no terminal
// (texto legível ou JSON, para integração com outras ferramentas/CI).
package report

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dev-rafaelsilva/secretscan/internal/patterns"
	"github.com/dev-rafaelsilva/secretscan/internal/scanner"
)

func Print(findings []scanner.Finding, missingGitignore []string, format string) error {
	if format == "json" {
		return printJSON(findings, missingGitignore)
	}
	printText(findings, missingGitignore)
	return nil
}

func printText(findings []scanner.Finding, missingGitignore []string) {
	if len(findings) == 0 {
		fmt.Println("✅ Nenhum segredo encontrado.")
	} else {
		for _, f := range findings {
			fmt.Printf("[%s] %s:%d → %s detectado\n", f.Severity, f.File, f.Line, f.Pattern)
		}
	}

	if len(missingGitignore) > 0 {
		fmt.Printf("\n⚠️  .gitignore: %d padrão(ões) sensível(is) ausente(s): %v\n", len(missingGitignore), missingGitignore)
	}

	high, medium := countBySeverity(findings)
	fmt.Printf("\n%d problema(s) encontrado(s) (%d alto, %d médio)\n", len(findings), high, medium)
}

func printJSON(findings []scanner.Finding, missingGitignore []string) error {
	out := struct {
		Findings         []scanner.Finding `json:"findings"`
		MissingGitignore []string          `json:"missing_gitignore"`
		Total            int               `json:"total"`
	}{findings, missingGitignore, len(findings)}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func countBySeverity(findings []scanner.Finding) (high, medium int) {
	for _, f := range findings {
		switch f.Severity {
		case patterns.High:
			high++
		case patterns.Medium:
			medium++
		}
	}
	return
}

func HasCritical(findings []scanner.Finding) bool {
	for _, f := range findings {
		if f.Severity == patterns.High {
			return true
		}
	}
	return false
}

// PrintHistory formata os achados do "secretscan history".
func PrintHistory(findings []scanner.HistoryFinding) {
	if len(findings) == 0 {
		fmt.Println("✅ Nenhum segredo encontrado no histórico de commits.")
		return
	}

	high := 0
	for _, f := range findings {
		fmt.Printf("[%s] commit %s — %s → %s detectado\n", f.Severity, f.Commit, f.File, f.Pattern)
		if f.Severity == patterns.High {
			high++
		}
	}
	fmt.Printf("\n%d problema(s) encontrado(s) no histórico (%d alto)\n", len(findings), high)
	fmt.Println("⚠️  Esses segredos ainda estão acessíveis via 'git log', mesmo se já foram removidos do código atual.")
	fmt.Println("   Rotacione as chaves encontradas — apagar o arquivo não é suficiente.")
}

// HasCriticalHistory retorna true se houver achado de severidade alta no histórico.
func HasCriticalHistory(findings []scanner.HistoryFinding) bool {
	for _, f := range findings {
		if f.Severity == patterns.High {
			return true
		}
	}
	return false
}