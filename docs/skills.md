# Skills — Padrão Skills.sh

## O que é Skills.sh

Skills.sh é um ecossistema aberto criado pela Vercel (2026) que define como agentes de IA executam ações reutilizáveis via linha de comando. É o "npm para agentes de IA".

Cada skill é uma pasta com:
- `SKILL.md` — metadata + instruções
- `run.sh` — script executável
- Arquivos auxiliares opcionais

---

## Formato SKILL.md

```yaml
---
name: web-search
description: "Busca informações na web"
version: "1.0.0"
author: slmpack
triggers: ["buscar", "pesquisar", "procurar", "google"]
needs_vram: false
needs_web: true
dependencies: ["curl", "jq"]
input: "query via stdin (texto)"
output: "JSON via stdout"
---
```

```markdown
# Web Search

Busca informações na web usando DuckDuckGo API.

## Uso
echo "últimas notícias IA" | bash run.sh

## Output
JSON array: [{title, url, snippet}]

## Exemplo
$ echo "cotação dólar hoje" | bash run.sh
[{"title": "Dólar hoje...", "url": "...", "snippet": "..."}]
```

---

## Skills Disponíveis

### web-search

Busca na web via DuckDuckGo.

```bash
echo "query" | bash skills/web-search/run.sh
```

### url-fetch

Baixa conteúdo de uma URL e extrai texto.

```bash
echo "https://example.com" | bash skills/url-fetch/run.sh
```

### file-read

Lê conteúdo de um arquivo.

```bash
echo "/path/to/file.txt" | bash skills/file-read/run.sh
```

### file-write

Escreve conteúdo em arquivo.

```bash
echo -e "/path/to/file.txt\nconteúdo aqui" | bash skills/file-write/run.sh
```

### rag-query

Consulta RAG (ChromaDB).

```bash
echo "pergunta sobre meus docs" | bash skills/rag-query/run.sh
```

### rag-ingest

Indexa documentos no RAG.

```bash
echo "/path/to/docs/" | bash skills/rag-ingest/run.sh
```

### code-execute

Executa código Python em sandbox.

```bash
echo "print('hello')" | bash skills/code-execute/run.sh
```

### pdf-read

Extrai texto de PDF.

```bash
echo "/path/to/file.pdf" | bash skills/pdf-read/run.sh
```

---

## Como o Agente Usa Skills

O flow engine (Go) executa skills via subprocess:

```go
// Pseudocode Go
func runSkill(skillName string, input string) (string, error) {
    cmd := exec.Command("bash", "skills/"+skillName+"/run.sh")
    cmd.Stdin = strings.NewReader(input)
    output, err := cmd.Output()
    return string(output), err
}
```

O agente (LLM) pode invocar skills via chamada de tool:

```
Usuário: "procure notícias sobre IA"
Router: detecta "procure" → flow "search"
Flow: executa skill web-search → pega resultados
Flow: passa resultados pro LLM (qwen3:4b) sintetizar
LLM: gera resposta baseada nos resultados
```

---

## Composição de Skills (Unix Pipes)

Skills são compostas via pipes Unix:

```bash
# Buscar + resumir
echo "notícias IA" | bash skills/web-search/run.sh | bash skills/summarize/run.sh

# Ler arquivo + extrair dados
cat documento.txt | bash skills/data-extract/run.sh

# RAG + código
echo "função que ordena" | bash skills/rag-query/run.sh | bash skills/code-generate/run.sh
```

Isso é o poder do padrão Unix aplicado a agentes de IA.

---

## Criar Nova Skill

1. Criar pasta em `skills/nome-da-skill/`
2. Criar `SKILL.md` com metadata
3. Criar `run.sh` executável
4. Adicionar triggers no `config/flows.toml`

```bash
mkdir -p skills/minha-skill
cat > skills/minha-skill/SKILL.md << 'EOF'
---
name: minha-skill
description: "Descrição da skill"
version: "1.0.0"
triggers: ["palavra1", "palavra2"]
---
# Minha Skill
Instruções aqui.
EOF

cat > skills/minha-skill/run.sh << 'EOF'
#!/bin/bash
INPUT=$(cat)
# processar...
echo "$INPUT"
EOF

chmod +x skills/minha-skill/run.sh
```
