# Arquitetura — Micro-LLM Modular

## Conceito Central

NÃO usar uma LLM grande que faz tudo. Em vez disso:

- **1 SLM router** (1-2B) decide qual especialista chamar
- **N especialistas** (modelos/pipelines) executam tarefas específicas
- **RAG** fornece contexto externo quando necessário
- **Tools** fazem a ponte com mundo real (web, filesystem, código)

A SLM router fica **sempre carregada na VRAM** (~0.5-1GB). Os especialistas carregam **sob demanda** e são descarregados quando não usados.

---

## Diagrama de Fluxo

```
Usuário → Input
         │
         ▼
┌─────────────────────────────┐
│  SLM Router (1.7B, ~50ms)  │  ← sempre na VRAM
│  "O que o usuário quer?"    │
└──────────┬──────────────────┘
           │
    ┌──────┴──────────────────────────────┐
    │          Classificação               │
    │  ┌──────┐ ┌──────┐ ┌──────┐ ┌─────┐ │
    │  │ CODAR│ │RAG   │ │WEB   │ │CHAT │ │
    │  │      │ │BUSCA │ │SEARCH│ │     │ │
    │  └──┬───┘ └──┬───┘ └──┬───┘ └──┬──┘ │
    └─────┼────────┼────────┼────────┼────┘
          │        │        │        │
          ▼        ▼        ▼        ▼
     ┌────────┐┌────────┐┌────────┐┌────────┐
     │DeepSeek││ RAG    ││SearxNG ││Qwen3   │
     │Coder   ││Pipeline││+Scraper││4B      │
     │(sob    ││(embed +││(sob    ││(sob    │
     │demanda)││query)  ││demanda)││demanda)│
     └────┬───┘└────┬───┘└────┬───┘└────┬───┘
          │         │         │         │
          └────┬────┴────┬────┴────┬────┘
               │         │         │
               ▼         ▼         ▼
          ┌────────────────────────────┐
          │    Response Synthesizer    │
          │  (formata, resume, entrega)│
          └────────────────────────────┘
```

---

## Regras de Memória (4GB VRAM)

### Estado 1: Ocioso (default)
```
SLM Router 1.7B (Q4):  ~0.8 GB
OS overhead:            ~0.5 GB
Buffer:                 ~2.7 GB livre
```

### Estado 2: Executando tarefa
```
SLM Router 1.7B:        ~0.8 GB
Especialista 4B (Q4):   ~2.5 GB
OS overhead:            ~0.5 GB
Buffer:                 ~0.2 GB  ← apertado mas funciona
```

### Estado 3: RAG ativo
```
SLM Router 1.7B:        ~0.8 GB
Embeddings (nomic):     ~0.3 GB
Query context:          ~0.3 GB
Buffer:                 ~2.6 GB livre  ← confortável
```

**Chave:** Nunca ter 2 especialistas 4B carregados ao mesmo tempo. O router descarrega um antes de carregar outro.

---

## Camadas do Sistema

### Camada 1 — Router (sempre ativo)
- SLM de 1-2B (Qwen3 1.7B ou Gemma 3 1B)
- Classifica intenção do usuário
- Decide qual especialista invocar
- Extrai parâmetros da query
- Latência: 50-200ms

### Camada 2 — Especialistas (sob demanda)
- Modelos especializados por domínio
- Carregam/descarregam dinamicamente
- Cada um tem sua própria configuração
- Ver `docs/especialistas.md`

### Camada 3 — RAG (memória externa)
- ChromaDB + embeddings
- Indexação de documentos locais
- Consulta semântica
- Ver `docs/rag.md`

### Camada 4 — Tools (ponte com mundo real)
- Web search
- Filesystem
- Python execution
- Browser automation
- Ver `docs/tools.md` (quando criado)

---

## Comparação: Monolítico vs Modular

| Aspecto | Monolítico (1 modelo grande) | Modular (router + especialistas) |
|---|---|---|
| VRAM | Sempre ~4GB+ | ~0.8GB ocioso, ~3.3GB ativo |
| Velocidade chat | 8-15 tok/s | 25-40 tok/s (router) + especialista |
| Qualidade código | Ruim se modelo genérico | Bom (modelo especializado) |
| Adicionar skill | Retreinar modelo | Adicionar novo especialista |
| Latência roteamento | N/A | +50-200ms (aceitável) |
| Flexibilidade | Baixa | Alta |

---

## Padrão Tiny-Critic (inspirado no paper)

Em vez de rodar o modelo grande para avaliar se o RAG retornou algo útil:

1. Router SLM avalia qualidade do retrieval (< 50ms)
2. Se bom → passa direto pro sintetizador
3. Se ruim → ativa fallback (web search, tool alternativa)

Isso evita gastar tokens do modelo grande em avaliação.

---

## Referências

- RouteLLM (LMSYS, ICLR 2025) — routing entre modelos com preference data
- vLLM Semantic Router (2026) — signal-driven decision routing
- Tiny-Critic RAG (arXiv:2603.00846) — SLM como gatekeeper
- Mixture of Agents (MoA) — colaboração entre múltiplos LLMs
- AgentForge (arXiv:2601.13383) — framework modular leve
