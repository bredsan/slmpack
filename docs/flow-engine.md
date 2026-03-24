# Flow Engine — Router Heurístico (sem modelo)

## Conceito

Em vez de usar uma SLM de 1.7B como router, usamos **código puro** para decidir. Regex, keywords, tamanho do input, padrões. Zero VRAM. Sub-milissegundo.

É um "n8n enxuto" em Python puro. Um flow engine que:
1. Recebe input do usuário
2. Analisa com regras heurísticas
3. Decide qual especialista chamar
4. Executa o flow via terminal

O LLM **só é chamado quando realmente precisa gerar texto**.

---

## Por que Heurístico > SLM Router?

| Aspecto | SLM Router (1.7B) | Heurístico (código) |
|---|---|---|
| VRAM | ~0.8GB | 0GB |
| Latência | 50-200ms | <1ms |
| Precisão | ~85-90% | ~80-95% (depende das regras) |
| Manutenção | Re-treinar modelo | Editar regras em Python |
| Custo | GPU necessária | CPU apenas |

Para **roteamento de intenção**, regras heurísticas são suficientes na maioria dos casos. O artigo da SemEval 2026 mostra que "rule-based classifiers provaram ser mais confiáveis que classificação LLM-based" para detecção de tipo de pergunta.

---

## Arquitetura do Flow Engine

```
Input do Usuário (terminal)
        │
        ▼
┌─────────────────────────┐
│   Heuristic Router      │  ← regex, keywords, tamanho
│   (código puro, <1ms)   │
└───────────┬─────────────┘
            │
            ▼
┌─────────────────────────┐
│   Flow Engine           │  ← executa o flow decidido
│   (n8n em Python)       │
└───────────┬─────────────┘
            │
     ┌──────┼──────┬──────┐
     ▼      ▼      ▼      ▼
  [CODER] [CHAT] [RAG] [WEB]  ← especialistas (sob demanda)
     │      │      │      │
     └──────┴──────┴──────┘
            │
            ▼
       Resposta final
```

---

## Regras Heurísticas (Router)

### Camada 1 — Detecção por Regex (instantâneo)

```python
ROUTES = {
    "codar": [
        r"\b(cod[eio]|função|function|def |class |import |script|programa|debug|erro|bug|implement|criar.*código)\b",
        r"```",                    # bloco de código
        r"\.(py|js|ts|html|css|go|rs|rb)\b",  # extensão de arquivo
    ],
    "resumir": [
        r"\b(resum[a-z]*|sintes[ei]|tl;?dr|em poucas palavras|chave|principais pontos)\b",
    ],
    "pesquisar": [
        r"\b(pesquis[a-z]*|busc[a-z]*|procur[a-z]*|google|notícia|últim[oa]s?|atualiz)\b",
        r"https?://",              # URL detectada
    ],
    "extrair": [
        r"\b(extra[ia][a-z]*|parse|json|csv|estrutur[a-z]*|lista de|tabela)\b",
    ],
    "raciocinar": [
        r"\b(prov[eia]|matemátic[a-z]|calcul[a-z]|resolv[a-z]|demonstr[a-z]|por que|explique.*por que)\b",
        r"[\d+\-*/=]{3,}",        # expressão matemática
    ],
    "ver": [
        r"\.(png|jpg|jpeg|gif|bmp|webp|svg)\b",  # imagem
        r"\b(imagem|foto|screenshot|tela|print)\b",
    ],
}
```

### Camada 2 — Detecção por Tamanho (instantâneo)

```python
COMPLEXITY_TIERS = {
    "simples":   lambda text: len(text) < 100 and "?" in text,
    "medio":     lambda text: 100 <= len(text) < 500,
    "complexo":  lambda text: len(text) >= 500,
}
```

### Camada 3 — Detecção por Palavras-chave (instantâneo)

```python
TOOL_TRIGGERS = {
    "needs_web": ["hoje", "agora", "atual", "último", "recente", "cotação", "preço"],
    "needs_rag": ["meu", "meus", "no meu", "nos meus", "documento", "arquivo", "nota"],
    "needs_file": ["abra", "leia", "salve", "escreva no", "crie o arquivo"],
}
```

---

## Flow Engine (mini-n8n)

O flow engine é um executor de steps em Python puro. Cada "flow" é uma lista de steps:

```python
from slmpack.flow import Flow, step

# Definir um flow
code_flow = Flow("coder", steps=[
    step.detect_intent,       # heurístico
    step.load_context,        # RAG se necessário
    step.call_specialist,     # carrega modelo, executa
    step.format_response,     # formata saída
])

# Executar
result = code_flow.run("escreva uma função Python que ordene uma lista")
```

### Step é uma função pura

```python
def detect_intent(ctx: Context) -> Context:
    """Regex + keywords → classifica intenção. <1ms."""
    text = ctx.input.lower()
    
    for intent, patterns in ROUTES.items():
        for pattern in patterns:
            if re.search(pattern, text):
                ctx.intent = intent
                ctx.confidence = 0.9
                return ctx
    
    ctx.intent = "conversar"  # fallback
    ctx.confidence = 0.5
    return ctx
```

### Flow é uma lista de steps

```python
class Flow:
    def __init__(self, name, steps):
        self.name = name
        self.steps = steps
    
    def run(self, user_input):
        ctx = Context(input=user_input)
        for step in self.steps:
            ctx = step(ctx)
            if ctx.done:
                break
        return ctx.output
```

### Flows pré-definidos

```python
FLOWS = {
    "codar": Flow("codar", [
        detect_intent,
        check_needs_rag,
        load_coder_model,      # qwen2.5-coder:3b
        generate_code,
        format_response,
    ]),
    "conversar": Flow("conversar", [
        detect_intent,
        check_needs_rag,
        load_chat_model,       # qwen3:4b
        generate_response,
        format_response,
    ]),
    "resumir": Flow("resumir", [
        detect_intent,
        load_summarizer_model, # qwen3:1.7b (leve)
        generate_summary,
        format_response,
    ]),
    "pesquisar": Flow("pesquisar", [
        detect_intent,
        execute_web_search,    # DuckDuckGo API
        load_chat_model,       # qwen3:4b (sintetiza)
        synthesize_results,
        format_response,
    ]),
}
```

---

## Interface Terminal

```bash
# Modo interativo
$ slmpack
slmpack> escreva uma função Python que calcule fatorial
[flow: codar] [model: qwen2.5-coder:3b] [vram: 3.3GB]
def factorial(n):
    if n <= 1:
        return 1
    return n * factorial(n - 1)

slmpack> resuma este texto: A inteligência artificial...
[flow: resumir] [model: qwen3:1.7b] [vram: 0.8GB]
A IA é um campo da computação focado em criar sistemas que simulam inteligência humana.

slmpack> qual a cotação do dólar hoje?
[flow: pesquisar] [needs_web: true] [model: qwen3:4b]
O dólar está cotado a R$ 5,15 hoje (24/03/2026).

slmpack> /status
Router: heuristic (<1ms)
Modelo ativo: qwen3:4b
VRAM: 3.3GB / 4.0GB
Flows executados: 42

slmpack> /flows
Flows disponíveis:
  codar      → qwen2.5-coder:3b
  conversar  → qwen3:4b
  resumir    → qwen3:1.7b
  pesquisar  → web + qwen3:4b
  raciocinar → qwen3:4b (thinking)
  extrair    → qwen3:1.7b
  ver        → gemma3:4b
```

---

## VRAM com Router Heurístico

```
Router heurístico:     0GB (código puro)
Especialista 4B:      ~2.5GB
OS overhead:           ~0.5GB
Buffer:                ~1.0GB  ← mais confortável que antes
─────────────────────────
Total:                 ~4.0GB
```

Comparado ao router SLM (0.8GB), ganhamos **0.8GB de buffer extra** — pode significar contexto maior ou modelo ligeiramente maior.

---

## Quando Heurística Falha

Se `confidence < 0.6` (nenhum regex bateu), o sistema:
1. Pergunta ao usuário: "Você quer que eu [codar/conversar/pesquisar/resumir]?"
2. Ou usa o modelo de chat (qwen3:4b) como fallback direto

Também pode ter um **modo debug** que mostra qual regra bateu:
```bash
slmpack> /debug escreva uma função
[regex match: codar pattern 1 → r"\b(função|function)\b"]
[confidence: 0.95]
[flow: codar]
```

---

## Referências

- vLLM Semantic Router — keyword routing com BM25, N-gram, regex (zero modelo)
- SemEval 2026 — "rule-based classifiers mais confiáveis que LLM-based"
- LogRocket — "Most teams discover they can handle 80% of routing needs with 5-10 simple rules"
- liteflow — workflow engine em Python puro com listas/dicts/sets
