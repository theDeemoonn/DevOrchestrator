package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"sync"
)

// availableTools содержит все доступные инструменты
var availableTools = map[string]Tool{
	"Node.js":            {"node", "Node.js", nil},
	"npm":                {"npm", "npm", nil},
	"Yarn":               {"yarn", "Yarn", nil},
	"Visual Studio Code": {"code", "Visual Studio Code", nil},
	"WebStorm":           {"webstorm", "WebStorm", nil},
	"Sublime Text":       {"sublime-text", "Sublime Text", nil},
	"OpenJDK":            {"java", "OpenJDK", nil},
	"Maven":              {"mvn", "Maven", nil},
	"Gradle":             {"gradle", "Gradle", nil},
	"IntelliJ IDEA":      {"intellij", "IntelliJ IDEA", nil},
	"Eclipse":            {"eclipse", "Eclipse", nil},
	"NetBeans":           {"netbeans", "NetBeans", nil},
	"Golang":             {"go", "Golang", nil},
	"Python 3":           {"python3", "Python 3", nil},
	"Pip":                {"pip3", "Pip", nil},
	"Virtualenv":         {"virtualenv", "Virtualenv", nil},
	"Git":                {"git", "Git", nil},
	"Docker":             {"docker", "Docker", nil},
	"Curl":               {"curl", "Curl", nil},
	"Zsh":                {"zsh", "Zsh", installOhMyZsh},
	"jq":                 {"jq", "jq", nil},
	"Postman":            {"postman", "Postman", nil},
	"Neovim":             {"nvim", "Neovim", installAstroNvim},
	"GoLand":             {"goland", "GoLand", nil},
	"PyCharm":            {"pycharm", "PyCharm", nil},
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "dev-installer",
		Short: "Утилита для установки инструментов разработчика",
		Long:  `Эта программа позволяет устанавливать, обновлять и удалять инструменты для различных стеков разработки.`,
		Run:   run,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Определяем ОС
	osType := detectOS()
	fmt.Printf("Обнаруженная ОС: %s\n", osType)

	// Шаг 1: Выбор действия
	action := selectAction()

	// Шаг 2: Выбор стека разработки
	stackStr := selectStack()
	stack := StringToStack(stackStr)

	// Шаг 3: Выбор инструментов
	var tools []string
	var ide []string

	if action == "Установить" {
		ide = selectIDE(stack)
	}
	tools = selectStackTools(stackStr)

	// Выполнение выбранного действия
	var err error
	switch action {
	case "Установить":
		err = performInstall(stack, ide, tools, osType)
	case "Обновить":
		err = performUpdate(tools, osType)
	case "Удалить":
		err = performUninstall(tools, osType)
	}

	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("%s: операция выполнена успешно.\n", action)
	}
}

func selectAction() string {
	prompt := promptui.Select{
		Label: "Выберите действие",
		Items: []string{"Установить", "Обновить", "Удалить"},
		Size:  3,
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Ошибка при выборе действия: %v", err)
	}

	return result
}

func performInstall(stack Stack, ide []string, tools []string, osType string) error {
	return installStack(stack, ide, tools, osType)
}

func performUpdate(tools []string, osType string) error {
	var wg sync.WaitGroup
	errorsCh := make(chan error, len(tools))

	for _, toolName := range tools {
		if tool, ok := availableTools[toolName]; ok {
			wg.Add(1)
			go func(tool Tool) {
				defer wg.Done()
				if err := tool.update(osType); err != nil {
					errorsCh <- fmt.Errorf("ошибка обновления %s: %v", tool.Description, err)
				}
			}(tool)
		}
	}

	wg.Wait()
	close(errorsCh)

	var errors []error
	for err := range errorsCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("произошли ошибки при обновлении: %v", errors)
	}
	return nil
}

func performUninstall(tools []string, osType string) error {
	var wg sync.WaitGroup
	errorsCh := make(chan error, len(tools))

	for _, toolName := range tools {
		if tool, ok := availableTools[toolName]; ok {
			wg.Add(1)
			go func(tool Tool) {
				defer wg.Done()
				if err := tool.uninstall(osType); err != nil {
					errorsCh <- fmt.Errorf("ошибка удаления %s: %v", tool.Description, err)
				}
			}(tool)
		}
	}

	wg.Wait()
	close(errorsCh)

	var errors []error
	for err := range errorsCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("произошли ошибки при удалении: %v", errors)
	}
	return nil
}

func selectStack() string {
	stacks := []string{
		"Frontend",
		"Java/Kotlin",
		"Golang",
		"Python",
		"Essential Tools",
	}

	prompt := promptui.Select{
		Label: "Выберите стек разработки",
		Items: stacks,
		Size:  len(stacks),
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Ошибка при выборе стека: %v", err)
	}

	return result
}

func selectStackTools(stack string) []string {
	toolsByStack := map[string][]string{
		"Frontend": {
			"Git",
			"Node.js",
			"npm",
			"Yarn",
			"Docker",
			"Curl",
			"Zsh",
			"jq",
			"Postman",
			"Neovim",
		},
		"Java/Kotlin": {
			"Git",
			"OpenJDK",
			"Maven",
			"Gradle",
			"Docker",
			"Curl",
			"Zsh",
			"jq",
			"Postman",
			"Neovim",
		},
		"Golang": {
			"Git",
			"Golang",
			"Docker",
			"Curl",
			"Zsh",
			"jq",
			"Postman",
			"Neovim",
		},
		"Python": {
			"Git",
			"Python 3",
			"Pip",
			"Virtualenv",
			"Docker",
			"Curl",
			"Zsh",
			"jq",
			"Postman",
			"Neovim",
		},
		"Essential Tools": {
			"Git",
			"Docker",
			"Curl",
			"Zsh",
			"jq",
			"Postman",
			"Neovim",
		},
	}

	prompt := promptui.Select{
		Label: "Выберите инструменты (Space для выбора, Enter для подтверждения)",
		Items: append([]string{"[Выбрать все]"}, toolsByStack[stack]...),
		Size:  len(toolsByStack[stack]) + 1,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\u25B6 {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "\u2714 {{ . | green }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Ошибка при выборе инструментов: %v", err)
	}

	if result == "[Выбрать все]" {
		return toolsByStack[stack]
	}

	return []string{result}
}

// Специальные функции установки
func installOhMyZsh() error {
	cmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка установки Oh My Zsh: %v, вывод: %s", err, string(output))
	}
	return nil
}

func installAstroNvim() error {
	cmd := exec.Command("git", "clone", "https://github.com/AstroNvim/AstroNvim", "~/.config/nvim")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка установки AstroNvim: %v, вывод: %s", err, string(output))
	}
	return nil
}
