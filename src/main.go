package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds the application configuration
type Config struct {
	Flows  []Flow  `toml:"flow"`
	Skills []Skill `toml:"skill"`
}

// Flow represents a workflow definition
type Flow struct {
	Name     string   `toml:"name"`
	Patterns []string `toml:"patterns"`
	Steps    []Step   `toml:"step"`
}

// Step represents a step in a flow
type Step struct {
	Skill string   `toml:"skill"`
	Args  []string `toml:"args,omitempty"`
	Model string   `toml:"model,omitempty"`
}

// Skill represents an executable skill
type Skill struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
	Desc string `toml:"desc"`
}

// App holds the application state
type App struct {
	Config    Config
	Skills    map[string]Skill
	Flows     []Flow
	Router    *Router
	OllamaURL string
}

// Router handles pattern matching for flows
type Router struct {
	patterns []*regexp.Regexp
	flows    []Flow
}

func NewRouter(flows []Flow) *Router {
	var patterns []*regexp.Regexp
	for _, f := range flows {
		for _, p := range f.Patterns {
			patterns = append(patterns, regexp.MustCompile(p))
		}
	}
	return &Router{patterns: patterns, flows: flows}
}

func (r *Router) Match(input string) []Flow {
	var matched []Flow
	for i, re := range r.patterns {
		if re.MatchString(input) {
			// Find which flow this pattern belongs to
			patternCount := 0
			for _, f := range r.flows {
				patternCount += len(f.Patterns)
				if i < patternCount {
					matched = append(matched, f)
					break
				}
			}
		}
	}
	return matched
}

func loadConfig() (Config, error) {
	var config Config
	_, err := toml.DecodeFile("config/slmpack.toml", &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func buildSkillMap(skills []Skill) map[string]Skill {
	skillMap := make(map[string]Skill)
	for _, s := range skills {
		skillMap[s.Name] = s
	}
	return skillMap
}

func executeStep(app *App, step Step, input string) (string, error) {
	skill, exists := app.Skills[step.Skill]
	if !exists {
		return "", fmt.Errorf("skill %s not found", step.Skill)
	}

	// Prepare command
	cmd := exec.Command(skill.Path)
	cmd.Args = append(cmd.Args, step.Args...)
	cmd.Args = append(cmd.Args, input)

	// Set environment for Ollama if needed
	if step.Model != "" {
		env := os.Environ()
		cmd.Env = append(env, "OLLAMA_MODEL="+step.Model)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("skill execution failed: %v", err)
	}
	return string(output), nil
}

func downloadFile(url string, filepath string) error {
	// Split the filepath to get the directory
	dir := filepath
	if lastSlash := strings.LastIndex(filepath, "/"); lastSlash >= 0 {
		dir = filepath[:lastSlash]
	}

	// Create the directory if it doesn't exist
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func setupWizard() {
	fmt.Println("Welcome to slmpack setup!")
	fmt.Println("This will help you install the default configuration and skills.")

	reader := bufio.NewReader(os.Stdin)

	// Ask if they want to download default config and skills
	fmt.Print("Download default configuration and skills from GitHub? [y/n]: ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		fmt.Println("Downloading files from GitHub...")
		baseURL := "https://raw.githubusercontent.com/bredsan/slmpack/master"

		files := []string{
			"config/slmpack.toml",
			"skills/code-execute/run.sh",
			"skills/web-search/run.sh",
			"skills/summarize/run.sh",
			"skills/vision-query/run.sh",
			"skills/chat/run.sh",
			"docs/status.md",
			"docs/flow-engine.md",
			"docs/skills.md",
			"docs/models.md",
			"docs/especialistas.md",
			"docs/rag.md",
			"docs/tools.md",
			"docs/router.md",
			"docs/fases.md",
			"docs/research.md",
			"docs/stack-tecnica.md",
			"docs/arquitetura.md",
		}

		for _, file := range files {
			url := baseURL + "/" + file
			fmt.Printf("Downloading %s... ", file)
			err := downloadFile(url, file)
			if err != nil {
				fmt.Printf("Failed: %v\n", err)
			} else {
				fmt.Println("OK")
			}
		}

		// Make skill scripts executable
		fmt.Println("Making skill scripts executable...")
		exec.Command("chmod", "+x", "skills/code-execute/run.sh").Run()
		exec.Command("chmod", "+x", "skills/web-search/run.sh").Run()
		exec.Command("chmod", "+x", "skills/summarize/run.sh").Run()
		exec.Command("chmod", "+x", "skills/vision-query/run.sh").Run()
		exec.Command("chmod", "+x", "skills/chat/run.sh").Run()
	}

	// Ask about Ollama models
	fmt.Print("\nDo you want to install Ollama models now? [y/n]: ")
	answer, _ = reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "y" || answer == "yes" {
		fmt.Println("Checking Ollama installation...")
		// Check if ollama command exists
		if _, err := exec.Command("ollama", "--version").Output(); err != nil {
			fmt.Println("Ollama is not installed or not in PATH. Please install Ollama first.")
			fmt.Println("Visit https://ollama.ai for installation instructions.")
		} else {
			// Check if Ollama server is running
			resp, err := http.Get("http://localhost:11434/api/version")
			if err != nil {
				fmt.Println("Ollama server is not running. Please start it with 'ollama serve' in another terminal.")
				fmt.Println("You can install models later by running: ollama pull <model>")
			} else {
				resp.Body.Close()
				fmt.Println("Ollama server is running.")
			}

			// List of models to choose from
			models := []string{
				"qwen3:4b",
				"qwen2.5-coder:3b",
				"qwen3:1.7b",
				"gemma3:4b",
				"nomic-embed-text",
			}

			fmt.Println("\nAvailable models:")
			for i, m := range models {
				fmt.Printf("  %d. %s\n", i+1, m)
			}
			fmt.Print("Enter the numbers of models to install (comma-separated, or 'all'): ")
			modelInput, _ := reader.ReadString('\n')
			modelInput = strings.TrimSpace(modelInput)

			var selectedModels []string
			if modelInput == "all" {
				selectedModels = models
			} else {
				parts := strings.Split(modelInput, ",")
				for _, part := range parts {
					numStr := strings.TrimSpace(part)
					if num, err := strconv.Atoi(numStr); err == nil && num > 0 && num <= len(models) {
						selectedModels = append(selectedModels, models[num-1])
					}
				}
			}

			if len(selectedModels) > 0 {
				fmt.Println("Installing selected models...")
				for _, model := range selectedModels {
					fmt.Printf("Pulling %s... ", model)
					if err := exec.Command("ollama", "pull", model).Run(); err != nil {
						fmt.Printf("Failed: %v\n", err)
					} else {
						fmt.Println("OK")
					}
				}
			}
		}
	}

	fmt.Println("\nSetup complete! You can now run slmpack.")
}

func main() {
	// Check if configuration exists
	if _, err := os.Stat("config/slmpack.toml"); os.IsNotExist(err) {
		setupWizard()
	}

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize app
	app := App{
		Config:    config,
		Skills:    buildSkillMap(config.Skills),
		Flows:     config.Flows,
		Router:    NewRouter(config.Flows),
		OllamaURL: "http://localhost:11434",
	}

	// Welcome message
	fmt.Println("slmpack - Local AI Pack for Limited Hardware")
	fmt.Println("Type '/help' for commands, '/exit' to quit")
	fmt.Println("")

	// REPL
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("slmpack> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		if input == "/exit" {
			break
		}

		if input == "/help" {
			printHelp()
			continue
		}

		if input == "/status" {
			printStatus(app)
			continue
		}

		if input == "/setup" {
			setupWizard()
			// After setup, reload config
			config, err = loadConfig()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reloading config: %v\n", err)
				os.Exit(1)
			}
			app.Config = config
			app.Skills = buildSkillMap(config.Skills)
			app.Flows = config.Flows
			app.Router = NewRouter(config.Flows)
			fmt.Println("Configuration reloaded. Continuing...")
			continue
		}

		// Route input to flows
		matchedFlows := app.Router.Match(input)
		if len(matchedFlows) == 0 {
			fmt.Println("[no matching flow]")
			continue
		}

		// Execute the first matching flow (in practice, could have priority)
		flow := matchedFlows[0]
		fmt.Printf("[flow:%s] ", flow.Name)

		// Execute flow steps
		currentInput := input
		for i, step := range flow.Steps {
			output, err := executeStep(&app, step, currentInput)
			if err != nil {
				fmt.Printf("\nError: %v\n", err)
				break
			}
			currentInput = output
			// For now, just print the final output
			if i == len(flow.Steps)-1 {
				fmt.Println(strings.TrimSpace(output))
			}
		}
	}
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  /help    - Show this help")
	fmt.Println("  /status  - Show system status")
	fmt.Println("  /setup   - Run setup wizard again")
	fmt.Println("  /exit    - Exit slmpack")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  escreva uma função Python que ordene uma lista")
	fmt.Println("  qual a cotação do dólar hoje?")
	fmt.Println("  resuma este texto: [cole seu texto aqui]")
}

func printStatus(app App) {
	fmt.Println("Router: heuristic (Go, <1ms)")
	fmt.Println("Model: default (via Ollama)")
	fmt.Println("VRAM: ~0.1GB / 4.0GB (Go core only)")
	fmt.Printf("Skills: %d disponíveis\n", len(app.Skills))
	fmt.Printf("Flows: %d configurados\n", len(app.Flows))
	// Check Ollama
	if _, err := exec.Command("curl", "-s", app.OllamaURL+"/api/version").Output(); err == nil {
		fmt.Println("Ollama: connected")
	} else {
		fmt.Println("Ollama: not connected (run 'ollama serve')")
	}

}
