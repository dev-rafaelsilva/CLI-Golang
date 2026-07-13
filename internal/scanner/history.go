package scanner

import (
	"fmt"
	"os"

	"github.com/dev-rafaelsilva/secretscan/internal/report"
	"github.com/dev-rafaelsilva/secretscan/internal/scanner"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Escaneia todo o histórico de commits do git em busca de segredos",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Escaneando histórico de commits (isso pode demorar em repos grandes)...")
		fmt.Println()

		findings, err := scanner.ScanHistory(".")
		if err != nil {
			return err
		}

		report.PrintHistory(findings)

		if report.HasCriticalHistory(findings) {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}