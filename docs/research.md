# Research — Referências e Papers

## Papers Acadêmicos

### Tiny-Critic RAG (2026)
- **Paper:** arXiv:2603.00846
- **Autores:** Yichao Wu et al.
- **Ideia:** SLM (Qwen-1.7B) com LoRA atua como gatekeeper para agentic RAG. Avalia se o retrieval é útil antes de passar para o modelo grande. Reduz latência em 10x comparado ao GPT-4o-mini.
- **Relevância:** Base da nossa arquitetura de routing.

### RouteLLM (ICLR 2025)
- **Paper:** arXiv:2406.18665
- **Autores:** LMSYS (Isaac Ong et al.)
- **Ideia:** Framework para routing entre LLMs usando preference data. Reduz custos em até 85% mantendo qualidade.
- **Relevância:** Prova que routing entre modelos é viável e eficiente.

### RouteMoA (2026)
- **Paper:** arXiv:2601.18130
- **Autores:** Jize Wang et al.
- **Ideia:** Mixture-of-Agents com routing dinâmico. Lightweight scorer filtra candidatos antes de inferência. Reduz custo em 89.8% e latência em 63.6%.
- **Relevância:** Confirma que pré-filtragem com modelo leve é eficaz.

### AgentForge (2026)
- **Paper:** arXiv:2601.13383
- **Autores:** A.A. Jafari, G. Anbarjafari
- **Ideia:** Framework modular leve para agents LLM. Skill abstraction com input/output contracts. YAML config. <100ms overhead de orquestração.
- **Relevância:** Inspiração para a estrutura modular do projeto.

### vLLM Semantic Router (2026)
- **Paper:** arXiv:2603.04444
- **Ideia:** Signal-driven decision routing para Mixture-of-Modality. 13 algoritmos de seleção. Routing sub-milissegundo.
- **Relevância:** Abordagem avançada de routing que podemos simplificar.

---

## Artigos e Guias

### Best Local LLM Models 2026 (SitePoint)
- Qwen3 7B lidera HumanEval na classe 7/8B
- Phi-4-mini (3.8B) é a única opção viável para 8GB machines
- Q4_K_M é o sweet spot de quantização

### Ollama VRAM Requirements (LocalLLM.in)
- 4GB VRAM → modelos 3-4B com Q4_K_M, contexto 4k tokens
- KV cache cresce linearmente com contexto
- VRAM é hard boundary, não soft limit

### Running LLMs on 4GB GPU (LinkedIn — Gabriele Monti)
- GTX 1050 Ti com 4GB VRAM: testou Gemma 3 2B, LLaMA2 7B Q4, DeepSeek R1 7B
- Modelos 7B funcionam mas são arriscados (estouram VRAM)
- Modelos 3-4B são o sweet spot para 4GB

### Model Router: Team of Small LLMs (Medium — Michael Hannecke)
- Padrão "recepção + especialistas sob demanda"
- Router 7B classifica em <300ms
- Especialistas carregam/descarregam dinamicamente
- Apple Silicon unified memory facilita, mas conceito se aplica a NVIDIA também

### Router-Based Agents (Towards AI)
- 3 níveis de routing: intra-agent (ReAct), inter-agent, orchestrator
- Semantic routing usa embeddings para dispatch
- LLM-based routing é mais caro mas mais preciso

---

## Ferramentas e Frameworks

### Ollama
- Runtime LLM local
- Suporte nativo a tool calling, embeddings, multimodal
- API REST simples
- https://ollama.com

### ChromaDB
- Banco de vetores para RAG
- Sem servidor separado (em processo)
- https://www.trychroma.com

### LangChain
- Framework de agentes
- Integração com Ollama
- https://langchain.com

### RouteLLM (framework)
- Framework open-source para routing entre LLMs
- Drop-in replacement para OpenAI client
- https://github.com/lm-sys/RouteLLM

---

## Modelos Específicos

### Qwen3 (Alibaba)
- 0.6B a 235B parâmetros
- Thinking mode nativo (ativa/desativa raciocínio)
- Tool use nativo
- Apache 2.0 license
- https://github.com/QwenLM/Qwen3

### Qwen2.5-Coder (Alibaba)
- Especializado em código
- 0.5B a 32B parâmetros
- Bom em Python, JavaScript, e outras linguagens

### Phi-4-mini (Microsoft)
- 3.8B parâmetros
- 128K contexto
- Forte em raciocínio e lógica
- MIT license

### Gemma 3 (Google)
- 270M a 27B parâmetros
- Multimodal (visão+texto)
- 128K contexto

---

## Links Úteis

- Ollama Library: https://ollama.com/library
- Hugging Face: https://huggingface.co
- LMSYS Chatbot Arena: https://chat.lmsys.org
- Papers With Code: https://paperswithcode.com
