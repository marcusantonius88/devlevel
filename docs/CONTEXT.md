# DevLevel - Contexto

## Ambiente de Execução

* A aplicação roda localmente na máquina do usuário
* Executada via terminal (CLI)
* Não requer container

---

## Linguagem e Execução

* Linguagem: Go
* Execução: `go run ./cmd`

---

## Integrações Externas

* GitHub REST API
* Autenticação via Personal Access Token

---

## Variáveis de Ambiente

* `GITHUB_TOKEN` deve ser fornecido

---

## Restrições

* Preferir biblioteca padrão do Go
* Evitar dependências externas desnecessárias
* Manter simplicidade e clareza

---

## Limitações do MVP

* Sem persistência
* Sem jobs em background
* Sem interface gráfica
