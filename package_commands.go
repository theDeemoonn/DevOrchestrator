// package_commands.go
package main

import "fmt"

// Обновляем команды для разных пакетных менеджеров
var packageManagerCommands = map[string]map[string]string{
	"apt": {
		"install":   "apt-get install -y",
		"uninstall": "apt-get remove -y",
		"update":    "apt-get upgrade -y",
	},
	"yum": {
		"install":   "yum install -y",
		"uninstall": "yum remove -y",
		"update":    "yum update -y",
	},
	"dnf": {
		"install":   "dnf install -y",
		"uninstall": "dnf remove -y",
		"update":    "dnf upgrade -y",
	},
	"pacman": {
		"install":   "-S --noconfirm",
		"uninstall": "-R --noconfirm",
		"update":    "-Syu --noconfirm",
	},
	"brew": {
		"install":   "install",
		"uninstall": "uninstall",
		"update":    "upgrade",
	},
	"choco": {
		"install":   "install -y",
		"uninstall": "uninstall -y",
		"update":    "upgrade -y",
	},
}

// Переопределяем установочные команды для специальных случаев
func getInstallCommand(program, osType string) string {
	switch osType {
	case "darwin":
		// Специальные случаи для macOS
		switch program {
		case "zsh":
			return "" // Oh My Zsh устанавливается через специальную функцию
		case "neovim":
			return "brew install neovim"
		default:
			return fmt.Sprintf("brew install %s", program)
		}
	default:
		return ""
	}
}

// Переопределяем команды обновления для специальных случаев
func getUpdateCommand(program, osType string) string {
	switch osType {
	case "darwin":
		// Специальные случаи для macOS
		switch program {
		case "zsh":
			return "upgrade_oh_my_zsh"
		case "neovim":
			return "brew upgrade neovim"
		default:
			return fmt.Sprintf("brew upgrade %s", program)
		}
	default:
		return ""
	}
}

// Переопределяем команды удаления для специальных случаев
func getUninstallCommand(program, osType string) string {
	switch osType {
	case "darwin":
		// Специальные случаи для macOS
		switch program {
		case "zsh":
			return "uninstall_oh_my_zsh"
		case "neovim":
			return "brew uninstall neovim"
		default:
			return fmt.Sprintf("brew uninstall %s", program)
		}
	default:
		return ""
	}
}
