package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "CLI-Golang",
	Short: "Escaneia repositórios em busca de segredos e configs sensíveis expostas",
	Long: `secretscan é uma CLI que analisa arquivos e histórico de commits
em busca de chaves de API, senhas e outros segredos que não deveriam
estar versionados no seu repositório.`,
	Run: func(cmd *cobra.Command, args []string) {
		runMenu()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}