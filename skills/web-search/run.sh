#!/bin/bash
# web-search skill - Search the web for information

INPUT="$*"

# For now, we'll just call Ollama with a prompt that simulates web search.
# In a real implementation, you might use Google Custom Search, DuckDuckGo, etc.

curl -s -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d "{
    \"model\": \"${OLLAMA_MODEL:-qwen3:4b}\",
    \"prompt\": \"Você é um assistente de busca na web. Forneça uma resposta concisa e atualizada para a pergunta: $INPUT\",
    \"stream\": false
  }" | jq -r '.response'