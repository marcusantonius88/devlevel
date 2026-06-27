# DevLevel - Especificação

## Visão Geral

DevLevel é uma ferramenta CLI que ajuda desenvolvedores a construírem consistência através de ofensivas diárias de commits no GitHub. Inspirada na mecânica de streak do Duolingo, transforma o hábito de commitar em uma experiência gamificada com streak tracking, daily goal, progressão de XP e ranks.

---

## Feature: Setup inicial

DADO QUE o usuário executa `devlevel setup` pela primeira vez
QUANDO o sistema solicitar o username do GitHub
ENTÃO deve salvar o username em `~/.devlevel/config.json`
E confirmar com a mensagem "✅ Configuration saved successfully"

DADO QUE o usuário executa `devlevel` sem ter feito setup
QUANDO o sistema não encontrar configuração local
ENTÃO deve exibir "❌ No GitHub username configured."
E orientar o usuário a executar `devlevel setup`

---

## Feature: Buscar Atividade

DADO QUE o username está configurado
QUANDO o sistema for executado
ENTÃO deve buscar os repos com atividade via `/users/{username}/events/public`
E para cada repo, buscar commits do autor via `/repos/{owner}/{repo}/commits`
E considerar apenas commits dos últimos 30 dias

DADO QUE a API retorna erro 403
QUANDO o sistema detectar rate limit
ENTÃO deve exibir mensagem explicativa sobre o limite de 60 requests/hora
E orientar o usuário a tentar novamente em alguns minutos

DADO QUE um ou mais repos não respondem dentro do timeout
QUANDO todos os repos falharem
ENTÃO deve exibir "Your streak is safe — please try again in a few minutes"
E não exibir dados zerados que possam confundir o usuário

DADO QUE apenas alguns repos falham por timeout
QUANDO parte dos dados for recuperada
ENTÃO deve exibir aviso de dados incompletos
E mostrar os stats com os dados disponíveis

---

## Feature: Cálculo de XP

DADO uma lista de commits novos (SHAs ainda não contabilizados)
QUANDO o sistema processar os commits
ENTÃO cada commit novo deve somar 10 XP ao total acumulado
E o XP total deve ser persistido em `~/.devlevel/state.json`
E o XP nunca deve diminuir entre execuções

Exemplo:

* 5 commits novos → +50 XP

---

## Feature: Cálculo de Nível (Craft Track)

DADO um valor total de XP acumulado
QUANDO o sistema calcular o nível
ENTÃO deve retornar o nível correspondente à tabela Craft Track:

| Level | Rank        | XP Mínimo |
|-------|-------------|-----------|
| 1     | Apprentice  | 0         |
| 2     | Craftsman   | 100       |
| 3     | Artisan     | 250       |
| 4     | Forger      | 500       |
| 5     | Blacksmith  | 750       |
| 6     | Grandmaster | 1000      |
| 7     | Sage        | 1500      |
| 8     | Oracle      | 2000      |
| 9     | Mythic      | 3000      |

---

## Feature: Cálculo de Streak

DADO os dias com atividade persistidos no state local
QUANDO o sistema calcular o streak
ENTÃO deve contar os dias consecutivos com pelo menos 1 commit, retroativamente a partir de hoje
E se não houver commit hoje, deve iniciar a contagem a partir de ontem (proteção de streak mid-day)
E deve usar o timezone local da máquina do usuário

Exemplo:

* Atividade em 3 dias consecutivos → Streak = 3
* Falhar um dia → streak reinicia do zero

---

## Feature: Daily Goal

DADO que os commits do dia foram processados
QUANDO o sistema verificar o objetivo diário
ENTÃO deve exibir "✅ Daily Goal: COMPLETE" se houver pelo menos 1 commit hoje (timezone local)
E deve exibir "⚠️ Daily Goal: PENDING — commit today to protect your streak" caso contrário

---

## Feature: Saída no CLI

DADO que os dados foram processados
QUANDO o sistema exibir as informações
ENTÃO a saída deve seguir o layout abaixo, com streak como elemento principal:

```
🚀 DevLevel
ℹ️  Using public GitHub API

🔥 CURRENT STREAK: N DAYS
✅ Daily Goal: COMPLETE

👤 User    : <username>
🏆 Level   : <level>
⚡ XP      : <xp>
🏅 Rank    : <rank>

📈 Progress to Level N
   ██████░░░░ N%
🎯 Next Level: N XP remaining

📊 Summary
   • Last 30 days: N commits
   • <mensagem motivacional contextual>
```

A mensagem motivacional deve se adaptar ao contexto:
* Streak < 7 dias → "Daily goal completed — see you tomorrow"
* Streak 7–29 dias → "Keep the momentum going"
* Streak 30+ dias → "Incredible consistency — keep it up"
* Daily goal pendente → "Commit today to protect your streak"

---

## Feature: Tratamento de Erros

DADO QUE o usuário não executou o setup
QUANDO o sistema iniciar
ENTÃO deve exibir mensagem orientando a executar `devlevel setup`

DADO QUE a API retorna rate limit (403)
QUANDO o sistema detectar o erro
ENTÃO deve exibir mensagem amigável explicando o limite e pedindo para tentar mais tarde

DADO QUE todos os repos atingem timeout
QUANDO nenhum dado for recuperado
ENTÃO deve informar que é um problema de comunicação e que o streak está seguro

---

## Fora de Escopo (MVP)

* Sem banco de dados externo
* Sem frontend
* Sem múltiplos usuários simultâneos
* Sem suporte a repositórios privados (limitação da API pública)
* Sem OAuth ou autenticação via token
