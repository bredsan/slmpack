# slmpack

> Pack de SLMs modulares para hardware limitado. Core em Go + Skills em Shell + Flows em TOML. Cabe em 4GB VRAM.

Sistema de IA local com **router heurístico (Go) + flow engine (Go) + skills (shell) + especialistas (Ollama)**. Zero overhead de roteamento. Toda VRAM para modelos.

**Hardware alvo:** GTX 1050 Ti (4GB VRAM) — ou qualquer GPU com 4GB+

---

## Stack

| Camada | Tecnologia | Função |
|---|---|---|
| Core | **Go** | Router heurístico + flow engine (binário único) |
| Skills | **Shell (.sh)** | Ações reutilizáveis (web, file, rag, code) |
| Config | **TOML** | Definição de flows, agents, skills |
| LLM | **Ollama** | Runtime dos modelos especializados |
| RAG | **ChromaDB + Python** | Memória externa (via skill) |

Python NÃO faz parte do core. Só aparece nas skills que chamam Ollama/ChromaDB.

---

## Arquitetura

```
Terminal (usuário)
      │
      ▼
┌─────────────────────────────┐
│  slmpack (binário Go)       │
│                             │
│  Router (regex, <1ms)       │  ← zero VRAM
│  Flow Engine (goroutines)   │  ← lê config TOML
│                             │
│  skills/web-search/run.sh   │  ← shell
│  skills/rag-query/run.sh    │  ← shell + python
│  skills/code-execute/run.sh │  ← shell
│                             │
│  → Ollama (localhost:11434) │  ← HTTP, modelos LLM
└─────────────────────────────┘
```

---

## Modelos (Ollama, sob demanda)

| Função | Modelo | VRAM | Tag |
|---|---|---|---|
| Chat | Qwen3 4B | ~2.5GB | `qwen3:4b` |
| Código | Qwen2.5-Coder 3B | ~2GB | `qwen2.5-coder:3b` |
| Resumo | Qwen3 1.7B | ~0.8GB | `qwen3:1.7b` |
| Visão | Gemma 3 4B | ~3GB | `gemma3:4b` |
| Embeddings | nomic-embed-text | ~0.3GB | `nomic-embed-text` |

---

## VRAM Budget (4GB)

| Estado | Uso | Livre |
|---|---|---|
| Ocioso (só Go, zero modelo) | ~0.1GB | ~3.9GB |
| Especialista 4B ativo | ~2.5GB | ~1.5GB |
| Especialista 4B + RAG | ~2.8GB | ~1.2GB |

Router usa 0GB VRAM. Go usa ~8MB RAM. Tudo para modelos.

---

## Velocidade Esperada (GTX 1050 Ti)

| Componente | Latência |
|---|---|
| Router (regex/Go) | <1ms |
| Flow dispatch | <1ms |
| Skill shell | ~10-50ms |
| Ollama 4B first token | ~1.5s |
| Ollama 1.7B first token | ~0.5s |
| RAG lookup | ~200ms |

---

## Terminal

```bash
$ slmpack
slmpack> escreva uma função Python que ordene uma lista
[flow:coder] [model:qwen2.5-coder:3b] [vram:2.5GB]
def sort_list(lst):
    return sorted(lst)

slmpack> qual a cotação do dólar hoje?
[flow:search] [skill:web-search] [model:qwen3:4b]
O dólar está a R$ 5,15 hoje (24/03/2026).

slmpack> /status
Router: heuristic (Go, <1ms)
Modelo: qwen3:4b
VRAM: 2.5GB / 4.0GB
Skills: 8 disponíveis
Flows: 5 configurados
```

---

## Documentação

| Arquivo | Conteúdo |
|---|---|
| [docs/stack-tecnica.md](docs/stack-tecnica.md) | Go + Shell + TOML — por quê |
| [docs/skills.md](docs/skills.md) | Padrão skills.sh, criar skills |
| [docs/flow-engine.md](docs/flow-engine.md) | Router heurístico + flow engine |
| [docs/arquitetura.md](docs/arquitetura.md) | Visão geral da arquitetura |
| [docs/modelos.md](docs/modelos.md) | Modelos por domínio |
| [docs/especialistas.md](docs/especialistas.md) | Pipelines especializadas |
| [docs/rag.md](docs/rag.md) | Sistema RAG |
| [docs/tools.md](docs/tools.md) | Web, filesystem, python exec |
| [docs/router.md](docs/router.md) | Router SLM (legado) |
| [docs/fases.md](docs/fases.md) | Roadmap de implementação |
| [docs/research.md](docs/research.md) | Papers, referências |
| [docs/status.md](docs/status.md) | Progresso atual |

---

## Referências

- [Bifrost](https://github.com/boundaryml/bifrost) — Go LLM gateway, 50x mais rápido que Python
- [Skills.sh](https://skills.sh) — ecossistema aberto de skills para agentes
- [Axe](https://github.com/jrswab/axe) — CLI Go para agentes single-purpose em TOML
- [vLLM Semantic Router](https://vllm-semantic-router.com) — keyword routing zero-overhead
