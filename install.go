// install.go
package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Установить инструменты для выбранного стека",
	Run:   installRun,
}

func installRun(cmd *cobra.Command, args []string) {
	osType := detectOS()
	fmt.Printf("Обнаруженная ОС: %s\n", osType)

	stackStr := selectStack()
	stack := StringToStack(stackStr)
	ide := selectIDE(stack)
	additionalTools := selectStackTools(stackStr)

	err := installStack(stack, ide, additionalTools, osType)
	if err != nil {
		fmt.Printf("Ошибка установки: %v\n", err)
	} else {
		fmt.Println("Установка завершена успешно.")
	}
}

func selectIDE(stack Stack) []string {
	var options []string

	switch stack {
	case FrontendStack:
		options = []string{"Visual Studio Code", "WebStorm", "Sublime Text"}
	case JavaKotlinStack:
		options = []string{"IntelliJ IDEA", "Eclipse", "NetBeans"}
	case GolangStack:
		options = []string{"Visual Studio Code", "GoLand", "Sublime Text"}
	case PythonStack:
		options = []string{"PyCharm", "Visual Studio Code", "Sublime Text"}
	default:
		options = []string{"Visual Studio Code", "Sublime Text"}
	}

	prompt := promptui.Select{
		Label: "Выберите IDE",
		Items: append([]string{"[Выбрать все]"}, options...),
		Size:  len(options) + 1,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "\u25B6 {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "\u2714 {{ . | green }}",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Ошибка при выборе IDE: %v", err)
	}

	if result == "[Выбрать все]" {
		return options
	}

	return []string{result}
}
