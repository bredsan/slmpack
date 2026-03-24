# Stack Técnica — Escolha de Linguagens

## Decisão: Go + Shell + Python (cada um no seu lugar)

| Camada | Linguagem | Motivo |
|---|---|---|
| **Core (router, flow engine)** | **Go** | 50x mais rápido que Python, ~11µs/request, goroutines, binário único |
| **Skills (ações reutilizáveis)** | **Shell (.sh)** | Padrão skills.sh da Vercel, zero overhead, composável com Unix |
| **Especialistas (LLM)** | **Python** | Ollama API, ecossistema ML, só onde precisa de modelo |
| **Config (flows, agents)** | **TOML** | Legível, tipado, como o projeto Axe |

---

## Por que Go para o Core?

### Benchmarks reais (2026)

| Métrica | Python (FastAPI) | Go (Fiber) | Diferença |
|---|---|---|---|
| Overhead por request | ~8ms | ~11µs | **700x mais rápido** |
| P99 latência @ 5k RPS | degrada (segundos) | 1.6-1.7s estável | **50x mais estável** |
| Memória | ~340MB | ~68MB | **5x menos RAM** |
| Goroutines vs threads | GIL limita | 10k goroutines @ 2KB cada | sem comparação |

### Para o slmpack especificamente

O flow engine precisa:
- Classificar intenção em <1ms (regex/keywords)
- Gerenciar múltiplos flows simultâneos
- Chamar subprocessos (Ollama, scripts shell)
- Streaming de responses
- Zero overhead quando ocioso

Go entrega tudo isso. Python não.

---

## Por que Shell para Skills?

O padrão **skills.sh** da Vercel (2026) mostra que skills de agentes devem ser:
- **Shell scripts** — executáveis diretamente
- **Componíveis** — pipe stdin/stdout
- **Documentadas em SKILL.md** — metadata + instruções
- **Sem dependências** — rodam em qualquer Unix

### Exemplo de Skill

```
skills/
├── web-search/
│   ├── SKILL.md          # metadata + instruções
│   └── run.sh            # execução real
├── code-review/
│   ├── SKILL.md
│   └── run.sh
├── file-summary/
│   ├── SKILL.md
│   └── run.sh
└── rag-query/
    ├── SKILL.md
    └── run.sh
```

### SKILL.md (formato skills.sh)

```yaml
---
name: web-search
description: "Busca informações na web via DuckDuckGo"
version: "1.0.0"
triggers: ["buscar", "pesquisar", "procurar", "google", "notícia"]
needs_vram: false
dependencies: ["curl", "jq"]
---
```

```markdown
# Web Search

Busca informações na web e retorna resultados estruturados.

## Uso
echo "query aqui" | bash skills/web-search/run.sh

## Output
JSON com [{title, url, snippet}]
```

### run.sh

```bash
#!/bin/bash
# Web Search Skill
# Input: query via stdin
# Output: JSON via stdout

QUERY=$(cat)
ENCODED=$(printf '%s' "$QUERY" | jq -sRr @uri)
curl -s "https://api.duckduckgo.com/?q=${ENCODED}&format=json" | \
  jq '[.RelatedTopics[:5][] | {title: .Text, url: .FirstURL}]'
```

Vantagens:
- Zero overhead (processo shell nativo)
- Componível com pipes Unix (`echo "query" | web-search | summarize`)
- O agente pode chamar via `bash skills/web-search/run.sh`
- Latência: microssegundos para iniciar

---

## Por que TOML para Config?

```toml
# config/flows.toml

[flow.coder]
description = "Gera código"
triggers = ["codar", "código", "function", "def ", "class ", "```"]
model = "qwen2.5-coder:3b"
vram = "2GB"
skills = ["file-read", "file-write"]

[flow.chat]
description = "Conversação geral"
triggers = []  # fallback quando nenhum outro bate
model = "qwen3:4b"
vram = "2.5GB"
skills = ["rag-query", "web-search"]

[flow.summarizer]
description = "Resumir textos"
triggers = ["resumir", "sintese", "tl;dr", "principais pontos"]
model = "qwen3:1.7b"
vram = "0.8GB"
skills = []

[flow.search]
description = "Busca na web"
triggers = ["pesquisar", "buscar", "procurar", "notícia", "hoje", "cotação"]
model = "qwen3:4b"
vram = "2.5GB"
skills = ["web-search", "url-fetch"]
needs_web = true
```

---

## Arquitetura Final

```
Terminal (usuário digita)
        │
        ▼
┌──────────────────────────────────┐
│  slmpack (binário Go)            │
│                                  │
│  ┌────────────────────────────┐  │
│  │ Heuristic Router           │  │  ← regex/keywords, <1ms
│  │ (Go, zero VRAM)            │  │
│  └──────────┬─────────────────┘  │
│             │                    │
│  ┌──────────▼─────────────────┐  │
│  │ Flow Engine                │  │  ← TOML config, executa steps
│  │ (Go, goroutines)           │  │
│  └──────────┬─────────────────┘  │
│             │                    │
│     ┌───────┼───────┐           │
│     ▼       ▼       ▼           │
│  [skills/] [ollama] [tools/]    │
│  shell      Python   shell      │
│  scripts    subprocess scripts   │
│             │                    │
└─────────────┼────────────────────┘
              │
              ▼
         Resposta
```

### Fluxo de execução

1. Usuário digita no terminal
2. Router Go (regex) classifica em <1ms
3. Flow Engine lê config TOML, decide qual flow
4. Flow executa steps:
   - Se precisa skill → `bash skills/X/run.sh`
   - Se precisa LLM → `curl localhost:11434/api/generate` (Ollama)
   - Se precisa tool → shell command
5. Resposta formatada ao usuário

---

## Comparação: Python vs Go para o Core

| Aspecto | Python | Go |
|---|---|---|
| Latência routing | ~50ms (import + regex) | ~0.01ms (compiled regex) |
| Memória base | ~30-50MB | ~5-8MB |
| Concorrência | GIL, asyncio complexo | goroutines nativas |
| Binário | precisa de Python instalado | binário único, sem dependências |
| Subprocessos | subprocess.Popen (lento) | os/exec (nativo, rápido) |
| Deploy | pip, venv, requirements | copiar binário |
| Startup | ~200ms (importar libs) | ~1ms |

---

## Onde Python Fica

Python **só** é usado nos especialistas que rodam via Ollama:
- Ollama já é um servidor (localhost:11434)
- Go chama via HTTP (curl ou net/http)
- Python não precisa rodar no core do slmpack
- O próprio Ollama é escrito em Go/C++

Então o stack real é:
- **Go**: core (router + flow engine + HTTP client)
- **Shell**: skills (scripts .sh)
- **TOML**: config (flows, agents, skills)
- **Ollama**: runtime LLM (já instalado, separado)

Python **não faz parte do slmpack**. O slmpack é 100% Go + Shell.

---

## Referências

- Bifrost (Go LLM gateway) — 50x mais rápido que LiteLLM (Python)
- Vercel Skills.sh — padrão aberto para skills de agentes via shell
- Axe (Go CLI) — agentes single-purpose definidos em TOML
- Benchmark Go vs Python vs Rust (2026) — Go ganha em I/O concorrente
- vLLM Semantic Router — keyword routing zero-overhead
