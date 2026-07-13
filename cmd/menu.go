package cmd

import (
	"fmt"
	"os"

	"github.com/dev-rafaelsilva/secretscan/internal/report"
	"github.com/dev-rafaelsilva/secretscan/internal/scanner"
	"github.com/manifoldco/promptui"
)

func runMenu() {
	prompt := promptui.Select{
		Label: "🔍 secretscan — o que você quer fazer?",
		Items: []string{
			"Escanear diretório atual",
			"Gerar .secretscan.yml",
			"Sair",
		},
	}

	_, escolha, err := prompt.Run()
	if err != nil {
		return
	}

	switch escolha {
	case "Escanear diretório atual":
		executarScan(".")
	case "Gerar .secretscan.yml":
		executarInit()
	case "Sair":
		fmt.Println("Até mais!")
	}
}

func executarScan(target string) {
	fmt.Printf("🔍 Escaneando %s...\n\n", target)

	findings, err := scanner.Scan(target)
	if err != nil {
		fmt.Fprintln(os.Stderr, "erro ao escanear:", err)
		return
	}
	missing := scanner.GitignoreMissing(target)
	report.Print(findings, missing, "text")
}

func executarInit() {
	if err := initCmd.RunE(initCmd, nil); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
	}
}