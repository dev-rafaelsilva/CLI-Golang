package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const defaultConfig = `# .secretscan.yml — configuração do secretscan
# Gerado automaticamente por "secretscan init"

# Diretórios que o scan nunca deve entrar
ignore_dirs:
  - .git
  - node_modules
  - vendor
  - dist
  - build
  - .venv

# Arquivos sensíveis que precisam estar no .gitignore
sensitive_files:
  - .env
  - "*.pem"
  - "*.key"
  - "*.p12"

# Achados que você já revisou e sabe que são falso positivo
# (ex: chave de exemplo no README). Formato: "arquivo:linha"
allowlist: []
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Gera um arquivo .secretscan.yml de configuração no diretório atual",
	RunE: func(cmd *cobra.Command, args []string) error {
		const path = ".secretscan.yml"

		if _, err := os.Stat(path); err == nil {
			fmt.Printf("⚠️  %s já existe. Nada foi alterado.\n", path)
			return nil
		}

		if err := os.WriteFile(path, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("erro ao criar %s: %w", path, err)
		}

		fmt.Printf("✅ %s criado com sucesso.\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}