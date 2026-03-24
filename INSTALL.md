# Instalação do slmpack

## Pré-requisitos

- GPU com pelo menos 4GB de VRAM (GTX 1050 Ti ou equivalente)
- [Ollama](https://ollama.ai) instalado e em execução
- Sistema operacional Windows 10/11 ou Linux

## Opções de Instalação

### Opção 1: Download do Release (Recomendado)

1. Acesse a página de releases: https://github.com/bredsan/slmpack/releases
2. Baixe o arquivo `slmpack.exe` para Windows ou o binário apropriado para seu sistema
3. Coloque o executável em uma pasta de sua escolha
4. Execute-o diretamente: `slmpack.exe`

### Opção 2: Instalação via Go

Se você tem o ambiente Go configurado:

```bash
# Clone o repositório
git clone https://github.com/bredsan/slmpack.git
cd slmpack

# Instale dependências
go mod tidy

# Compile o executável
go build -o slmpack.exe src/main.go

# Execute
slmpack.exe
```

### Opção 3: Instalação via Script (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/bredsan/slmpack/main/install.sh | bash
```

## Configuração Inicial

Após instalar, você precisa configurar os modelos Ollama:

```bash
# Baixe os modelos necessários (isso pode levar algum tempo)
ollama pull qwen3:4b
ollama pull qwen2.5-coder:3b
ollama pull qwen3:1.7b
ollama pull gemma3:4b
ollama pull nomic-embed-text

# Inicie o Ollama em outro terminal
ollama serve
```

## Primeiro Uso

```bash
slmpack.exe
```

Digite suas solicitações no prompt. Exemplos:
- `escreva uma função Python que ordene uma lista`
- `qual a cotação do dólar hoje?`
- `/status` para ver informações do sistema
- `/help` para ver comandos disponíveis
- `/exit` para sair

## Solução de Problemas

### Ollama não está rodando
Certifique-se de que o Ollama está em execução:
```bash
ollama serve
```

### Erro de conexão com Ollama
Verifique se o Ollama está acessível em http://localhost:11434

### Modelo não encontrado
Certifique-se de que baixou o modelo correto:
```bash
ollama pull <nome-do-modelo>
```

### Performance baixa
O slmpack foi otimizado para GPUs com 4GB de VRAM. Para melhor desempenho:
- Feche aplicativos que consomem VRAM
- Verifique se a GPU está sendo utilizada (não integrada)
- Certifique-se de que drivers estão atualizados