# Modelos — Seleção por Domínio

## Hardware: GTX 1050 Ti — 4GB VRAM

Todos os modelos devem ser Q4_K_M ou menor. Modelos de 4B são o teto prático.

---

## Modelos Disponíveis via Ollama

### Router (sempre carregado)

| Modelo | Params | VRAM | Função | Tag |
|---|---|---|---|---|
| **Qwen3 1.7B** | 1.7B | ~0.8GB | Router principal, thinking mode | `qwen3:1.7b` |
| Gemma 3 1B | 1B | ~0.5GB | Alternativa ultra-leve | `gemma3:1b` |
| SmolLM2 1.7B | 1.7B | ~1.0GB | Boa qualidade geral | `smollm2:1.7b` |

**Recomendado:** `qwen3:1.7b` — tem thinking mode, tool use nativo, e cabe com folga.

### Especialistas (carregam sob demanda)

#### Código

| Modelo | Params | VRAM | Força | Tag |
|---|---|---|---|---|
| **Qwen2.5-Coder 3B** | 3B | ~2GB | Código geral, Python, JS | `qwen2.5-coder:3b` |
| DeepSeek-Coder 1.3B | 1.3B | ~0.8GB | Ultra rápido, código básico | `deepseek-coder:1.3b` |
| CodeLlama 7B Q3 | 7B | ~4GB | Fronteira, mais poderoso | `codellama:7b-code` |

#### Chat / Conversação

| Modelo | Params | VRAM | Força | Tag |
|---|---|---|---|---|
| **Qwen3 4B** | 4B | ~2.5GB | Melhor geral, thinking mode | `qwen3:4b` |
| Phi-4-mini | 3.8B | ~3GB | Lógica, raciocínio | `phi4-mini` |
| Gemma 3 4B | 4B | ~3GB | Multimodal (visão+texto) | `gemma3:4b` |

#### Raciocínio / Matemática

| Modelo | Params | VRAM | Força | Tag |
|---|---|---|---|---|
| **Qwen3 4B** | 4B | ~2.5GB | Thinking mode ativado | `qwen3:4b` |
| DeepSeek-R1 1.5B | 1.5B | ~1GB | Raciocínio leve | `deepseek-r1:1.5b` |

#### Resumo / Extração

| Modelo | Params | VRAM | Força | Tag |
|---|---|---|---|---|
| **Qwen3 1.7B** | 1.7B | ~0.8GB | Rápido, bom para resumos | `qwen3:1.7b` |
| Gemma 3 1B | 1B | ~0.5GB | Mais rápido ainda | `gemma3:1b` |

#### Classificação / Routing

| Modelo | Params | VRAM | Força | Tag |
|---|---|---|---|---|
| **Qwen3 1.7B** | 1.7B | ~0.8GB | Classificação de intenção | `qwen3:1.7b` |

#### Embeddings (RAG)

| Modelo | Params | VRAM | Função | Tag |
|---|---|---|---|---|
| **nomic-embed-text** | ~270MB | ~0.3GB | Embeddings para RAG | `nomic-embed-text` |
| all-MiniLM-L6-v2 | ~80MB | ~0.1GB | Alternativa mais leve | via sentence-transformers |

---

## Tabela de Decisão do Router

O router usa esta lógica para decidir qual especialista carregar:

```
Input do usuário
       │
       ▼
  [Qwen3 1.7B] ← classifica intenção
       │
       ├── "codar/escrever código/debug" → DeepSeek-Coder 1.3B ou Qwen2.5-Coder 3B
       ├── "explicar/conversar/perguntar" → Qwen3 4B
       ├── "resumir/extrair/chave" → Qwen3 1.7B (reusa o router)
       ├── "raciocínio/provar/matemática" → Qwen3 4B (thinking mode)
       ├── "pesquisar/buscar" → Tool: web search + RAG
       ├── "ver/ler imagem" → Gemma 3 4B (multimodal)
       └── "não sei" → Qwen3 4B (fallback geral)
```

---

## Configuração dos Modelos

### Router — Modelfile

```dockerfile
FROM qwen3:1.7b

PARAMETER num_ctx 1024
PARAMETER num_gpu 99
PARAMETER temperature 0.1
PARAMETER repeat_penalty 1.0

SYSTEM """Você é um classificador de intenção. Dada a mensagem do usuário, responda APENAS com JSON:
{"intencao": "codar|conversar|resumir|raciocinar|pesquisar|ver", "confianca": 0.0-1.0, "params": {}}"""
```

### Coder — Modelfile

```dockerfile
FROM qwen2.5-coder:3b

PARAMETER num_ctx 2048
PARAMETER num_gpu 99
PARAMETER temperature 0.3
PARAMETER repeat_penalty 1.1

SYSTEM """Você é um programador especializado. Escreva código limpo, eficiente e com nomes claros. Sempre inclua o que o código faz em 1 linha antes do código."""
```

### Chat — Modelfile

```dockerfile
FROM qwen3:4b

PARAMETER num_ctx 2048
PARAMETER num_gpu 99
PARAMETER temperature 0.7
PARAMETER repeat_penalty 1.1

SYSTEM """Você é um assistente útil e direto. Responda em português brasileiro. Seja conciso mas completo."""
```

---

## Comandos de Instalação

```bash
# Router
ollama pull qwen3:1.7b
ollama create micro-router -f Modelfile.router

# Especialistas
ollama pull qwen2.5-coder:3b     # código
ollama pull qwen3:4b              # chat geral
ollama pull deepseek-r1:1.5b      # raciocínio leve

# Embeddings
ollama pull nomic-embed-text

# Multimodal (opcional, se tiver VRAM sobrando)
ollama pull gemma3:4b
```
