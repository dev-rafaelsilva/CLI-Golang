/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dev-rafaelsilva/secretscan/internal/report"
	"github.com/dev-rafaelsilva/secretscan/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	strict bool
	format string
)

var runCmd = &cobra.Command{
	Use:   "run [diretório]",
	Short: "Escaneia um diretório em busca de segredos expostos",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		if format != "json" {
			fmt.Printf("🔍 Escaneando %s...\n\n", target)
		}

		findings, err := scanner.Scan(target)
		if err != nil {
			return fmt.Errorf("erro ao escanear: %w", err)
		}
		missing := scanner.GitignoreMissing(target)

		if err := report.Print(findings, missing, format); err != nil {
			return err
		}

		if report.HasCritical(findings) || (strict && len(findings) > 0) {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVar(&strict, "strict", false, "falha (exit 1) mesmo com achados de baixo risco")
	runCmd.Flags().StringVar(&format, "format", "text", "formato de saída: text | json")
}