# Especialistas — Modelos e Pipelines Especializados

## Conceito

Cada "especialista" é um modelo ou pipeline que resolve **um tipo de tarefa muito bem**. O router decide qual usar. O especialista executa e devolve o resultado.

Nenhum especialista fica carregado permanentemente. Todos carregam **sob demanda** e são descarregados após uso.

---

## Lista de Especialistas

### 1. Coder (`coder`)

**Função:** Escrever, debugar, explicar código

**Modelo:** `qwen2.5-coder:3b` (~2GB VRAM)

**Configuração:**
- temperature: 0.3 (precisão sobre criatividade)
- num_ctx: 2048

**Input esperado:**
```json
{"linguagem": "python", "descricao": "função que ordena lista"}
```

**Output:** Código + comentário de 1 linha

**Fallback:** Se modelo 3B não couber, usar `deepseek-coder:1.3b` (~0.8GB)

---

### 2. Chat (`chat`)

**Função:** Conversação geral, perguntas, explicações

**Modelo:** `qwen3:4b` (~2.5GB VRAM)

**Configuração:**
- temperature: 0.7
- num_ctx: 2048

**Input esperado:**
```json
{"mensagem": "explique o que é RAG"}
```

**Output:** Resposta em linguagem natural

---

### 3. Reasoner (`reasoner`)

**Função:** Matemática, lógica, provas, raciocínio complexo

**Modelo:** `qwen3:4b` com thinking mode ON (~2.5GB VRAM)

**Configuração:**
- temperature: 0.2
- num_ctx: 2048

**Input esperado:**
```json
{"problema": "prove que a raiz de 2 é irracional"}
```

**Output:** Raciocínio passo a passo + conclusão

**Nota:** Usa o mesmo modelo do chat, mas com thinking mode ativado e temperatura mais baixa.

---

### 4. Summarizer (`summarizer`)

**Função:** Resumir textos, extrair pontos-chave

**Modelo:** `qwen3:1.7b` (reusa o router, ~0.8GB VRAM)

**Configuração:**
- temperature: 0.3
- num_ctx: 4096 (contexto maior para textos longos)

**Input esperado:**
```json
{"texto": "...", "max_palavras": 100, "formato": "bullet_points"}
```

**Output:** Resumo formatado

---

### 5. Extractor (`extractor`)

**Função:** Extrair dados estruturados de texto

**Modelo:** `qwen3:1.7b` (reusa o router, ~0.8GB VRAM)

**Configuração:**
- temperature: 0.1
- num_ctx: 2048

**Input esperado:**
```json
{"texto": "...", "schema": {"nome": "string", "idade": "int"}}
```

**Output:** JSON estruturado

---

### 6. Vision (`vision`)

**Função:** Analisar imagens, screenshots

**Modelo:** `gemma3:4b` (~3GB VRAM)

**Configuração:**
- temperature: 0.5
- num_ctx: 2048

**Input esperado:**
```json
{"imagem": "path/to/image.png", "pergunta": "o que tem nesta imagem?"}
```

**Output:** Descrição/análise da imagem

---

### 7. Web Search Tool (`web_tool`)

**Função:** Buscar informações na web (não é modelo, é tool)

**Implementação:** DuckDuckGo API ou SearxNG

**Input esperado:**
```json
{"query": "últimas notícias sobre IA"}
```

**Output:** Lista de resultados com título, URL, snippet

---

## Pipeline de Execução

```
Router devolve JSON
       │
       ▼
Orchestrator lê "especialista"
       │
       ├── Se "coder" → carrega qwen2.5-coder:3b → executa → descarrega
       ├── Se "chat" → carrega qwen3:4b → executa → descarrega
       ├── Se "reasoner" → carrega qwen3:4b (thinking) → executa → descarrega
       ├── Se "summarizer" → usa qwen3:1.7b (já carregado) → executa
       ├── Se "extractor" → usa qwen3:1.7b (já carregado) → executa
       ├── Se "vision" → carrega gemma3:4b → executa → descarrega
       └── Se "web_tool" → chama API de busca → retorna
       │
       ▼
Response Synthesizer formata saída
```

---

## Gerenciamento de VRAM

### Sequência típica de carga/descarga

```
t=0ms    [VRAM: 0.8GB] Router carregado
t=50ms   [VRAM: 0.8GB] Router classifica → "codar"
t=60ms   [VRAM: 3.3GB] Carrega Coder 3B
t=1500ms [VRAM: 3.3GB] Coder gera código
t=2000ms [VRAM: 0.8GB] Descarrega Coder, Router volta ao controle
t=2050ms [VRAM: 0.8GB] Resposta entregue ao usuário
```

### Regra de Ouro
> Se o especialista 4B + router 1.7B > 4GB, o sistema usa CPU offload para o que não couber na VRAM. Ollama faz isso automaticamente.

### Configuração Ollama para Offload

```bash
# Forçar tudo na GPU quando possível
OLLAMA_GPU_LAYERS=99 ollama serve

# Ou no Modelfile
PARAMETER num_gpu 99
```

---

## Adicionar Novo Especialista

1. Escolher modelo que caiba na VRAM (≤3B para especialista + 1.7B router)
2. Criar Modelfile com system prompt especializado
3. Adicionar entrada na tabela de decisão do router
4. Testar latência e qualidade

Exemplo — adicionar especialista "tradutor":
```bash
ollama pull qwen2.5:3b  # modelo com bom multilíngue
```
```dockerfile
FROM qwen2.5:3b
PARAMETER num_ctx 2048
PARAMETER temperature 0.3
SYSTEM "Você é um tradutor profissional. Traduza fielmente."
```
Adicionar ao router:
```
- traduzir: traduzir texto entre idiomas → translator
```
