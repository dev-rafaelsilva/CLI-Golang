package cmd

import (
	"fmt"
	"os"

	"github.com/dev-rafaelsilva/secretscan/internal/report"
	"github.com/dev-rafaelsilva/secretscan/internal/scanner"
	"github.com/manifoldco/promptui"
)

// runMenu abre o menu interativo quando "secretscan" é chamado sem
// nenhum subcomando (nem "run", "init" etc).
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
		// Ctrl+C ou Esc cai aqui — sai calado, sem stacktrace feio.
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

// executarScan roda a mesma lógica do "secretscan run ." direto,
// sem passar pelo parsing de flags do Cobra.
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
	// reaproveita o RunE do initCmd já existente
	if err := initCmd.RunE(initCmd, nil); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
	}
}