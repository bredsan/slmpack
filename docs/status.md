# Status — Progresso do Projeto

## Data: 24/03/2026

### Concluído
- [x] Pesquisa de modelos para 4GB VRAM (Qwen3, Phi-4-mini, Gemma 3)
- [x] Pesquisa de arquiteturas modulares (RouteLLM, Tiny-Critic RAG, MoA)
- [x] Pesquisa de skills.sh (Vercel, padrão aberto para agentes)
- [x] Pesquisa Go vs Python para flow engine (Go 50x mais rápido)
- [x] Definição da stack: Go (core) + Shell (skills) + TOML (config)
- [x] Definição do nome: slmpack
- [x] Documentação completa em .md

### Decisões Técnicas

| Decisão | Escolha | Motivo |
|---|---|---|
| Core language | Go | 50x mais rápido que Python, ~11µs/request, binário único |
| Router | Heurístico (regex/Go) | Zero VRAM, <1ms, sem modelo |
| Skills | Shell scripts (.sh) | Padrão skills.sh, zero overhead, composável |
| Config | TOML | Legível, tipado, padrão Axe |
| Chat | Qwen3 4B | Melhor geral, thinking mode |
| Coder | Qwen2.5-Coder 3B | Especializado em código |
| RAG DB | ChromaDB | Simples, sem servidor |
| Embeddings | nomic-embed-text | Bom, leve, via Ollama |
| Runtime LLM | Ollama | Simples, API pronta, GPU auto-detect |
| Quantização | Q4_K_M | Melhor equilíbrio qualidade/VRAM |

### Stack Final

```
Go (core) → Shell (skills) → Ollama (LLM)
     ↓           ↓                ↓
  <1ms        ~10-50ms        ~1-2s
  0VRAM       0VRAM          2-3GB VRAM
```

### Pendente
- [ ] Instalar Go
- [ ] Instalar Ollama
- [ ] Pull especialistas (qwen3:4b, qwen2.5-coder:3b, qwen3:1.7b, nomic-embed-text)
- [ ] Criar estrutura Go do projeto (go.mod, main.go)
- [ ] Implementar heuristic router em Go
- [ ] Implementar flow engine em Go
- [ ] Criar skills shell (web-search, file-read, rag-query, etc.)
- [ ] Criar config TOML (flows.toml, agents.toml)
- [ ] Implementar terminal interface
- [ ] Testar fluxo completo
