#!/bin/bash
# code-execute skill - Generate and optionally execute code

# Get input from command line arguments
INPUT="$*"

# Default to just generating code, not executing
EXECUTE=false

# Check if we should execute (for now, we'll just generate)
if [[ "$1" == "--execute" ]]; then
    EXECUTE=true
    shift
    INPUT="$*"
fi

# Call Ollama to generate code
curl -s -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d "{
    \"model\": \"${OLLAMA_MODEL:-qwen2.5-coder:3b}\",
    \"prompt\": \"Você é um assistente especializado em geração de código. Gere apenas o código solicitado, sem explicações adicionais, a menos que explicitamente pedido. Entrada: $INPUT\",
    \"stream\": false
  }" | jq -r '.response'