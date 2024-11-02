package main

import (
	"fmt"
	"log"
)

// Stack представляет тип стека технологий
type Stack string

const (
	FrontendStack   Stack = "Frontend"
	JavaKotlinStack Stack = "Java/Kotlin"
	GolangStack     Stack = "Golang"
	PythonStack     Stack = "Python"
	EssentialStack  Stack = "Essential Tools"
)

// Tool представляет инструмент разработчика
type Tool struct {
	Command     string
	Description string
	InstallFunc func() error
}

// install устанавливает инструмент
func (t Tool) install(osType string) error {
	if !isInstalled(t.Command, osType) {
		if t.InstallFunc != nil {
			return t.InstallFunc()
		}
		return executeCommand(osType, "install", t.Command)
	}
	fmt.Printf("%s уже установлен.\n", t.Description)
	return nil
}

// update обновляет инструмент
func (t Tool) update(osType string) error {
	if isInstalled(t.Command, osType) {
		log.Printf("Обновление %s...\n", t.Description)
		return executeCommand(osType, "update", t.Command)
	}
	fmt.Printf("%s не установлен.\n", t.Description)
	return nil
}

// uninstall удаляет инструмент
func (t Tool) uninstall(osType string) error {
	if isInstalled(t.Command, osType) {
		log.Printf("Удаление %s...\n", t.Description)
		return executeCommand(osType, "uninstall", t.Command)
	}
	fmt.Printf("%s не установлен.\n", t.Description)
	return nil
}

// StringToStack конвертирует строку в тип Stack
func StringToStack(s string) Stack {
	switch s {
	case "Frontend":
		return FrontendStack
	case "Java/Kotlin":
		return JavaKotlinStack
	case "Golang":
		return GolangStack
	case "Python":
		return PythonStack
	case "Essential Tools":
		return EssentialStack
	default:
		return EssentialStack
	}
}
