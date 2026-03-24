package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

func main() {
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
