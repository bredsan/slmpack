# Fases — Roadmap de Implementação

## Fase 0 — Setup Base
**Tempo estimado:** 1-2 horas

- [ ] Instalar Ollama
- [ ] Pull dos modelos:
  - `ollama pull qwen3:1.7b` (router)
  - `ollama pull qwen3:4b` (chat)
  - `ollama pull nomic-embed-text` (embeddings)
- [ ] Testar cada modelo isoladamente
- [ ] Verificar VRAM usage com `nvidia-smi`
- [ ] Criar Modelfiles de configuração

**Entregável:** Modelos rodando localmente, chat básico funcional.

---

## Fase 1 — Router
**Tempo estimado:** 2-4 horas

- [ ] Implementar função `route(input) → dict`
- [ ] Definir schema JSON de saída do router
- [ ] Testar classificação com 10+ queries diferentes
- [ ] Ajustar prompt do sistema para consistência
- [ ] Implementar fallback quando confiança < 0.6
- [ ] Medir latência do router (alvo: < 200ms)

**Entregável:** Router classificando intenções corretamente.

---

## Fase 2 — Orchestrator + Especialistas
**Tempo estimado:** 4-6 horas

- [ ] Implementar orchestrator que lê JSON do router
- [ ] Implementar carga/descarga de modelos via Ollama API
- [ ] Implementar specialistas:
  - [ ] Coder (qwen2.5-coder:3b)
  - [ ] Chat (qwen3:4b)
  - [ ] Summarizer (qwen3:1.7b, reusa router)
- [ ] Implementar Response Synthesizer (formatação final)
- [ ] Testar fluxo completo: input → router → especialista → output
- [ ] Monitorar VRAM durante execução

**Entregável:** Sistema modular funcionando com 3 especialistas.

---

## Fase 3 — RAG
**Tempo estimado:** 3-5 horas

- [ ] Instalar ChromaDB (`pip install chromadb`)
- [ ] Implementar pipeline de ingesta:
  - [ ] Leitura de arquivos
  - [ ] Chunking
  - [ ] Embedding
  - [ ] Armazenamento no ChromaDB
- [ ] Implementar consulta RAG
- [ ] Integrar RAG com orchestrator
- [ ] Indexar documentos de exemplo
- [ ] Testar: pergunta sobre documentos locais

**Entregável:** Sistema consultando documentos locais antes de responder.

---

## Fase 4 — Tools
**Tempo estimado:** 3-5 horas

- [ ] Web search (DuckDuckGo API ou SearxNG)
- [ ] Filesystem (leitura/escrita de arquivos)
- [ ] Python exec (sandbox de execução)
- [ ] Integrar tools com router (needs_web, etc.)
- [ ] Implementar Tiny-Critic pattern:
  - [ ] Router avalia se RAG retornou resultado útil
  - [ ] Se não, ativa web search como fallback

**Entregável:** Sistema com acesso a web, arquivos e execução de código.

---

## Fase 5 — Refinamento
**Tempo estimado:** contínuo

- [ ] Fine-tuning LoRA do router (opcional, melhora classificação)
- [ ] Cache de respostas frequentes
- [ ] Métricas de performance (tok/s, latência, acurácia)
- [ ] Interface web simples (opcional)
- [ ] Logs e debugging
- [ ] Documentação do uso diário

---

## Fase 6 — Avançado (futuro)

- [ ] Multi-step planning (agent que decompõe tarefas complexas)
- [ ] Memória de conversa (histórico persistente)
- [ ] Tool chaining (resultado de uma tool vira input de outra)
- [ ] Especialistas customizados via fine-tuning
- [ ] Avaliação automática de qualidade das respostas

---

## Métricas de Sucesso

| Métrica | Alvo |
|---|---|
| Latência do router | < 200ms |
| Acurácia de classificação | > 85% |
| Tok/s (especialista 4B) | > 8 |
| Tok/s (router 1.7B) | > 25 |
| VRAM pico | < 3.8GB |
| Tempo de carga de especialista | < 3s |
| Qualidade RAG (recall@3) | > 70% |
