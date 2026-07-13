// Package patterns contém as assinaturas (regex) usadas para detectar
// segredos conhecidos em código-fonte.
package patterns

import "regexp"

type Severity string

const (
	High   Severity = "ALTO"
	Medium Severity = "MEDIO"
	Low    Severity = "BAIXO"
)

type Pattern struct {
	Name     string
	Regex    *regexp.Regexp
	Severity Severity
}

// Default retorna o conjunto padrão de padrões conhecidos.
func Default() []Pattern {
	return []Pattern{
		{
			Name:     "AWS Access Key",
			Regex:    regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
			Severity: High,
		},
		{
			Name:     "AWS Secret Key",
			Regex:    regexp.MustCompile(`(?i)aws_secret_access_key\s*[:=]\s*['"]?[A-Za-z0-9/+=]{40}['"]?`),
			Severity: High,
		},
		{
			Name:     "Stripe API Key",
			Regex:    regexp.MustCompile(`sk_live_[0-9a-zA-Z]{24,}`),
			Severity: High,
		},
		{
			Name:     "GitHub Token",
			Regex:    regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36,}`),
			Severity: High,
		},
		{
			Name:     "JWT Token",
			Regex:    regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`),
			Severity: Medium,
		},
		{
			Name:     "RSA Private Key",
			Regex:    regexp.MustCompile(`-----BEGIN (RSA )?PRIVATE KEY-----`),
			Severity: High,
		},
		{
			Name:     "Senha hardcoded em variável",
			Regex:    regexp.MustCompile(`(?i)(password|senha|passwd)\s*[:=]\s*['"][^'"]{4,}['"]`),
			Severity: Medium,
		},
		{
			Name:     "Generic API Key",
			Regex:    regexp.MustCompile(`(?i)(api[_-]?key)\s*[:=]\s*['"][A-Za-z0-9_\-]{16,}['"]`),
			Severity: Medium,
		},
	}
}