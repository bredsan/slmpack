# Router — Sistema de Roteamento

## O que é

O Router é uma SLM (Small Language Model) que fica **sempre carregada na VRAM** e decide qual especialista invocar para cada input do usuário.

Não é uma LLM "inteligente". É um **classificador rápido** que:
1. Lê a mensagem do usuário
2. Classifica a intenção
3. Extrai parâmetros relevantes
4. Retorna JSON com a decisão

Latência alvo: **50-200ms**

---

## Modelo Recomendado

**Qwen3 1.7B (Q4_K_M)** — ~0.8GB VRAM

Por quê:
- Tem thinking mode (pode desativar para ser mais rápido)
- Suporta tool use nativo
- Bom em classificação mesmo sendo pequeno
- Cabe com folga nos 4GB alongside qualquer especialista

Alternativa: Gemma 3 1B (~0.5GB) se precisar de mais espaço para especialistas.

---

## Schema de Saída

O router SEMPRE responde em JSON:

```json
{
  "intencao": "codar",
  "confianca": 0.92,
  "especialista": "coder",
  "modelo": "qwen2.5-coder:3b",
  "params": {
    "linguagem": "python",
    "descricao": "função que ordena lista"
  },
  "needs_rag": false,
  "needs_web": false
}
```

### Campos

| Campo | Tipo | Descrição |
|---|---|---|
| `intencao` | string | Categoria da tarefa |
| `confianca` | float | 0.0 a 1.0, confiança na classificação |
| `especialista` | string | Nome do especialista a invocar |
| `modelo` | string | Tag Ollama do modelo |
| `params` | object | Parâmetros extraídos da query |
| `needs_rag` | bool | Se precisa consultar RAG antes |
| `needs_web` | bool | Se precisa buscar na web |

---

## Intenções Suportadas

| Intenção | Especialista | Modelo | Quando |
|---|---|---|---|
| `codar` | coder | qwen2.5-coder:3b | Escrever, debugar, explicar código |
| `conversar` | chat | qwen3:4b | Conversa geral, perguntas |
| `resumir` | router | qwen3:1.7b | Resumir textos (reusa o router) |
| `raciocinar` | reasoner | qwen3:4b | Matemática, lógica, provas |
| `pesquisar` | web_tool | tool | Buscar na web |
| `ver` | vision | gemma3:4b | Analisar imagens |
| `extrair` | router | qwen3:1.7b | Extrair dados estruturados |
| `planejar` | chat | qwen3:4b | Criar planos, estratégias |

---

## Implementação

### Prompt do Router

```
Você é um classificador de intenção. Analise a mensagem do usuário e responda APENAS com JSON.

Intenções disponíveis:
- codar: escrever, modificar, debugar, ou explicar código
- conversar: perguntas gerais, bate-papo
- resumir: resumir ou condensar texto
- raciocinar: matemática, lógica, provas, problemas complexos
- pesquisar: buscar informações na web
- ver: analisar imagens, screenshots
- extrair: extrair dados estruturados de texto
- planejar: criar planos, organizar tarefas

Formato de resposta:
{"intencao": "...", "confianca": 0.0, "params": {}}

Mensagem do usuário: {input}
```

### Exemplo de Execução

```python
import ollama
import json

def route(input_text: str) -> dict:
    response = ollama.chat(
        model='qwen3:1.7b',
        messages=[{
            'role': 'system',
            'content': ROUTER_SYSTEM_PROMPT
        }, {
            'role': 'user',
            'content': input_text
        }],
        options={'temperature': 0.1}
    )
    
    return json.loads(response['message']['content'])

# Exemplo
decisao = route("escreva uma função Python que ordene uma lista")
# → {"intencao": "codar", "confianca": 0.95, "params": {"linguagem": "python"}}
```

---

## Fallback

Se `confianca < 0.6`, o router pode:
1. Perguntar ao usuário qual especialista usar
2. Usar o modelo de chat (qwen3:4b) como fallback
3. Combinar: RAG + chat genérico

---

## Otimização

### Thinking Mode OFF
Para classificação, desativar thinking mode do Qwen3:
```
/think off
```
Isso reduz latência de ~200ms para ~80ms.

### Cache de Classificações
Se o mesmo tipo de input aparece várias vezes, cachear a decisão:
```
"como fazer um for em python" → codar (cached)
```

### Batch de Decisões
Se múltiplos inputs chegarem juntos, classificar todos em uma única chamada.
