# DevLevel - Contexto

## Ambiente de Execução

* A aplicação roda localmente na máquina do usuário
* Executada via terminal (CLI)
* Não requer container

---

## Linguagem e Execução

* Linguagem: Go 1.21+
* Execução diária: `devlevel` (binário compilado)
* Execução a partir do fonte: `go run ./cmd`

---

## Integrações Externas

* GitHub REST API v3 (endpoints públicos — sem autenticação)
* Endpoints utilizados:
  * `GET /users/{username}/events/public` — descobre repos com atividade recente
  * `GET /repos/{owner}/{repo}/commits?author={username}&since={date}` — busca commits por repo

---

## Configuração do Usuário

* Username do GitHub é salvo localmente via `devlevel setup`
* Arquivo de configuração: `~/.devlevel/config.json`
* Sem variáveis de ambiente obrigatórias
* Sem tokens ou autenticação

---

## Persistência Local

* Estado de progressão salvo em `~/.devlevel/state.json`
* Armazena:
  * XP total acumulado (nunca diminui)
  * SHAs de commits já contabilizados (evita dupla contagem)
  * Dias com atividade (permite calcular streaks além da janela da API)

---

## Restrições

* Usar apenas biblioteca padrão do Go (zero dependências externas)
* Manter simplicidade e clareza de código
* Arquitetura hexagonal: domínio isolado de infraestrutura e apresentação

---

## Limitações conhecidas

* Apenas repositórios públicos são contabilizados (limitação da API pública)
* API pública limita 60 requests/hora por IP — uso intensivo pode atingir o limite
* Janela de eventos da API cobre aproximadamente os últimos 30 dias
