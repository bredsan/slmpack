# RAG — Retrieval Augmented Generation

## O que é

RAG permite que o sistema consulte documentos locais antes de responder. Em vez de depender apenas do conhecimento treinado no modelo, ele busca informação relevante na sua base de dados.

---

## Arquitetura

```
Documento (PDF, TXT, MD, código)
       │
       ▼
┌──────────────┐
│   Chunking   │ ← divide em pedaços de ~512 tokens
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Embedding   │ ← nomic-embed-text converte em vetor
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  ChromaDB    │ ← armazena vetores + texto original
└──────────────┘


Consulta do usuário
       │
       ▼
┌──────────────┐
│  Embedding   │ ← mesma função de embedding
│  da query    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Similarity  │ ← busca vetorial (cosine similarity)
│  Search      │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Top-K Docs  │ ← retorna os N chunks mais relevantes
└──────┬───────┘
       │
       ▼
  Contexto injetado no prompt do especialista
```

---

## Componentes

### Embedding Model

**Modelo:** `nomic-embed-text` via Ollama

- Tamanho: ~270MB
- VRAM: ~0.3GB (carrega sob demanda, não precisa ficar permanente)
- Dimensão: 768
- Suporte multilíngue

**Alternativa mais leve:** `all-MiniLM-L6-v2` via sentence-transformers (~80MB)

### Vector Store

**Ferramenta:** ChromaDB

Por quê:
- Mais simples de configurar
- Não precisa de servidor separado (roda em processo)
- Bom para uso pessoal/protótipo
- Persiste em disco

Alternativa: FAISS (mais rápido para grandes volumes, mas mais complexo)

### Chunking

- Tamanho: 512 tokens por chunk
- Overlap: 20% (102 tokens)
- Separadores: parágrafos, quebras de linha

---

## Setup

```bash
# Instalar dependências
pip install chromadb ollama

# O modelo de embedding já deve estar instalado
ollama pull nomic-embed-text
```

---

## Fluxo de Ingestão

```python
import chromadb
import ollama

# 1. Criar/conectar ao banco
client = chromadb.PersistentClient(path="./data/chroma")
collection = client.get_or_create_collection("docs")

# 2. Ler documento
with open("documento.txt", "r") as f:
    texto = f.read()

# 3. Chunking simples
chunks = []
chunk_size = 512  # caracteres (simplificado)
overlap = 100
for i in range(0, len(texto), chunk_size - overlap):
    chunks.append(texto[i:i + chunk_size])

# 4. Gerar embeddings
embeddings = []
for chunk in chunks:
    response = ollama.embeddings(model='nomic-embed-text', prompt=chunk)
    embeddings.append(response['embedding'])

# 5. Armazenar
collection.add(
    embeddings=embeddings,
    documents=chunks,
    ids=[f"chunk_{i}" for i in range(len(chunks))]
)
```

---

## Fluxo de Consulta

```python
def query_rag(pergunta: str, top_k: int = 3) -> list[str]:
    # 1. Embedding da pergunta
    response = ollama.embeddings(model='nomic-embed-text', prompt=pergunta)
    query_embedding = response['embedding']
    
    # 2. Busca vetorial
    results = collection.query(
        query_embeddings=[query_embedding],
        n_results=top_k
    )
    
    # 3. Retorna chunks relevantes
    return results['documents'][0]
```

---

## Integração com Router

Quando o router detecta que a query pode se beneficiar de contexto externo:

```json
{
  "intencao": "conversar",
  "needs_rag": true,
  "params": {"topico": "meus projetos"}
}
```

O orchestrator:
1. Consulta RAG primeiro
2. Pega os chunks relevantes
3. Injeta no prompt do especialista
4. Especialista responde com contexto

---

## Fontes de Dados

Colocar documentos em `data/documents/`:
- PDFs
- Arquivos .txt, .md
- Código fonte
- Notas (Obsidian, etc.)
- Logs

O sistema indexa automaticamente ou sob demanda.

---

## Otimização para 4GB VRAM

O modelo de embedding (nomic-embed-text) NÃO precisa ficar carregado permanentemente:
- Carrega para ingesta/indexação
- Carrega para consulta
- Descarrega após uso
- VRAM liberada para especialistas

Custo de carregar/descarregar: ~200ms. Aceitável.
