# Tools — Ponte com o Mundo Real

## Conceito

Tools não são modelos. São **funções executáveis** que o sistema chama quando precisa interagir com o mundo externo: web, filesystem, terminal, APIs.

O router decide *se* precisa de uma tool. O orchestrator chama a tool. O resultado volta para o especialista processar.

---

## Tools Disponíveis

### 1. Web Search

**Função:** Buscar informações na web

**Implementação:**
- DuckDuckGo API (sem API key, mais simples)
- SearxNG (self-hosted, mais privado)

**Input:** `{"query": "o que é transformer em ML"}`

**Output:** Lista de resultados `{title, url, snippet}`

**Quando o router ativa:** `needs_web: true`

---

### 2. Filesystem

**Função:** Ler e escrever arquivos locais

**Operações:**
- `read_file(path)` — lê conteúdo de arquivo
- `write_file(path, content)` — escreve arquivo
- `list_dir(path)` — lista diretório
- `search_files(pattern)` — busca por nome/conteúdo

**Quando o router ativa:** intenção `codar` com params indicando arquivo

---

### 3. Python Exec

**Função:** Executar código Python em sandbox

**Implementação:** Subprocess com timeout e restrições

**Limitações:**
- Sem acesso a rede (por padrão)
- Sem acesso a arquivos fora de workspace
- Timeout de 30s
- Sem imports perigosos (os.system, subprocess direto)

**Quando o router ativa:** intenção `codar` com `executar: true`

---

### 4. URL Fetcher

**Função:** Baixar conteúdo de uma URL

**Input:** `{"url": "https://example.com/article"}`

**Output:** Texto extraído da página (sem HTML)

**Quando o router ativa:** URL detectada na mensagem do usuário

---

### 5. PDF Reader

**Função:** Extrair texto de arquivos PDF

**Dependência:** PyPDF2 ou pdfplumber

**Quando o router ativa:** intenção `resumir` ou `extrair` com arquivo .pdf

---

## Integração com Router

O router inclui flags `needs_rag` e `needs_web` no JSON de saída:

```json
{
  "intencao": "conversar",
  "confianca": 0.85,
  "especialista": "chat",
  "modelo": "qwen3:4b",
  "needs_rag": true,
  "needs_web": false,
  "params": {}
}
```

O orchestrator executa na ordem:
1. Se `needs_rag` → consulta RAG → injeta contexto
2. Se `needs_web` → chama web search → injeta resultados
3. Chama especialista com contexto enriquecido
4. Retorna resposta

---

## Padrão Tiny-Critic para Tools

Em vez de sempre chamar a tool, o router pode avaliar se é necessário:

```
Usuário: "qual a capital da França?"
Router: needs_rag=false, needs_web=false  ← sabe responder direto

Usuário: "qual a cotação do dólar hoje?"
Router: needs_web=true  ← precisa buscar, é dado dinâmico

Usuário: "o que está no meu documento sobre X?"
Router: needs_rag=true  ← precisa consultar base local
```

Isso evita chamadas desnecessárias de tools, reduzindo latência.
