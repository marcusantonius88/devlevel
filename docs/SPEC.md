# DevLevel - Especificação

## Visão Geral

DevLevel é uma ferramenta CLI que gamifica a atividade de desenvolvimento com base nos commits do GitHub.

---

## Feature: Buscar Atividade

DADO QUE um token válido do GitHub foi fornecido
QUANDO o sistema for executado
ENTÃO deve buscar os commits do usuário autenticado
E considerar apenas commits dos últimos 7 dias

---

## Feature: Cálculo de XP

DADO uma lista de commits
QUANDO o sistema calcular XP
ENTÃO cada commit deve somar 10 XP

Exemplo:

* 5 commits → 50 XP

---

## Feature: Cálculo de Nível

DADO um valor total de XP
QUANDO o sistema calcular o nível
ENTÃO deve retornar:

* Level 1 → XP < 100
* Level 2 → XP >= 100
* Level 3 → XP >= 250
* Level 4 → XP >= 500

---

## Feature: Cálculo de Streak

DADO commits agrupados por dia
QUANDO o sistema calcular o streak
ENTÃO o streak deve ser a quantidade de dias consecutivos com pelo menos 1 commit

Exemplo:

* Atividade em 3 dias consecutivos → Streak = 3
* Falhar um dia → streak reinicia

---

## Feature: Daily Goal

DADO que os commits do dia foram processados
QUANDO o sistema verificar o objetivo diário
ENTÃO deve exibir COMPLETE se houver pelo menos 1 commit hoje
E deve exibir PENDING caso contrário

---

## Feature: Saída no CLI

DADO que os dados foram processados
QUANDO o sistema exibir as informações
ENTÃO a saída deve seguir o layout:

```
🚀 DevLevel

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
   • Recent activity: N commits
   • <motivational message>
```

---

## Feature: Tratamento de Erros

DADO QUE o token do GitHub não foi informado
QUANDO o sistema iniciar
ENTÃO deve exibir uma mensagem de erro
E encerrar a execução

---

## Fora de Escopo (MVP)

* Sem banco de dados
* Sem frontend
* Sem cache
* Sem múltiplos usuários
