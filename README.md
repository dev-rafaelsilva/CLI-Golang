# secretscan

CLI em Go que escaneia repositórios em busca de segredos expostos — chaves de API, senhas, tokens — tanto no código atual quanto no **histórico de commits**, onde a maioria das ferramentas não olha.

🔍 **Escaneando...**

| Severidade | Arquivo | Tipo de segredo |
|---|---|---|
| 🔴 Alto | `.env:2` | Stripe API Key |
| 🟡 Médio | `docker-compose.yml:4` | Generic API Key |
| 🔴 Alto | `src/config.js:2` | AWS Access Key |

⚠️ **`.gitignore`**: 4 padrões sensíveis ausentes — `.env`, `*.pem`, `*.key`, `*.p12`

**3 problemas encontrados** (2 alto, 1 médio)

## Por que isso importa

Apagar um arquivo com uma chave exposta **não remove ela do histórico do Git**. Qualquer pessoa com acesso ao repositório ainda pode encontrar essa chave rodando `git log -p`. O `secretscan history` existe exatamente pra pegar esse ponto cego — ele escaneia todos os commits, não só o estado atual dos arquivos.

## Instalação

Você tem duas opções. Escolha uma.

### Opção 1 — `go install` (recomendado, se você já tem Go instalado)

```bash
go install github.com/dev-rafaelsilva/secretscan@latest
```

Isso baixa, compila e instala o binário automaticamente na pasta de binários do Go (`$GOPATH/bin` ou `%USERPROFILE%\go\bin` no Windows). Depois disso, o comando `secretscan` funciona em qualquer lugar do terminal, sem precisar de `./` na frente.

**Se depois de instalar o terminal disser "comando não reconhecido"**, é porque essa pasta de binários do Go ainda não está no PATH do seu sistema:

- **Windows**: rode `go env GOPATH` para ver o caminho (geralmente `C:\Users\SEU_USUARIO\go`). Adicione `C:\Users\SEU_USUARIO\go\bin` às variáveis de ambiente (Painel de Controle → Editar variáveis de ambiente do sistema → Path → Novo). Feche e reabra o terminal depois.
- **Mac/Linux**: adicione `export PATH=$PATH:$(go env GOPATH)/bin` ao seu `~/.zshrc` ou `~/.bashrc`, depois rode `source ~/.zshrc` (ou reabra o terminal).

### Opção 2 — Clonar e compilar manualmente

```bash
git clone https://github.com/dev-rafaelsilva/secretscan.git
cd secretscan
go build -o secretscan .
```

No **Windows**, o Go sempre gera o binário com extensão `.exe`, então use:

```bash
go build -o secretscan.exe .
```

Com essa opção, o binário fica só na pasta atual (não no PATH), então você precisa rodar com `./` (Mac/Linux) ou `.\` (Windows) na frente — veja a seção de uso abaixo.

## Uso

Se você instalou com `go install` (Opção 1), roda direto:

```bash
secretscan run .            # escaneia o diretório atual
secretscan history           # escaneia todo o histórico de commits
secretscan init              # gera um .secretscan.yml de configuração
secretscan                   # abre o menu interativo
```

Se você compilou manualmente (Opção 2), precisa indicar que o binário está na pasta atual:

```bash
# Mac/Linux
./secretscan run .

# Windows
.\secretscan.exe run .
```

### Modo estrito e CI/CD

```bash
secretscan run . --strict
```

Retorna exit code `1` se encontrar qualquer achado — inclusive severidade baixa. Isso permite travar um pipeline automaticamente:

```yaml
# .github/workflows/secretscan.yml
- name: Scan for secrets
  run: |
    go install github.com/dev-rafaelsilva/secretscan@latest
    secretscan run . --strict
```

## O que ele detecta

- Chaves de acesso e secret keys da AWS
- Chaves de API da Stripe
- Tokens de acesso do GitHub
- Tokens JWT
- Chaves privadas RSA
- Senhas hardcoded em variáveis
- Chaves de API genéricas
- Arquivos sensíveis (`.env`, `*.pem`, `*.key`, `*.p12`) ausentes do `.gitignore`

## Problemas comuns

| Erro | Causa | Solução |
|---|---|---|
| `repository not found` no `git clone` | URL digitada errada ou nome do repositório mudou | Confira a URL exata na página do GitHub |
| `'secretscan' não é reconhecido como um comando` (Windows) | Faltou o `.\` na frente, ou faltou compilar com `.exe` | Use `.\secretscan.exe run .` |
| `command not found` (Mac/Linux) após `go install` | A pasta de binários do Go não está no PATH | Veja a seção de instalação acima |
| `git não encontrado ou não é um repositório git` ao rodar `secretscan history` | O comando `history` precisa ser rodado dentro de uma pasta com `git init` já feito | Rode `git init` na pasta antes de usar `history` |

## Status

🚧 Em desenvolvimento ativo. Feedback, issues e contribuições são muito bem-vindos.
