// utils.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Карта специальных команд установки
var specialInstallCommands = map[string]map[string][]string{
	"windows": {
		"code": {
			"powershell -Command \"choco install vscode -y\"",
		},
		"sublime-text": {
			"powershell -Command \"choco install sublimetext3 -y\"",
		},
		"goland": {
			"powershell -Command \"choco install goland -y\"",
		},
		"webstorm": {
			"powershell -Command \"choco install webstorm -y\"",
		},
		"pycharm": {
			"powershell -Command \"choco install pycharm-community -y\"",
		},
		"intellij-idea": {
			"powershell -Command \"choco install intellijidea-community -y\"",
		},
		"docker": {
			"powershell -Command \"choco install docker-desktop -y\"",
		},
	},
	"linux": {
		"code": {
			"wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg",
			"sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg",
			"sudo sh -c 'echo \"deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main\" > /etc/apt/sources.list.d/vscode.list'",
			"rm -f packages.microsoft.gpg",
			"sudo apt update",
			"sudo apt install -y code",
		},
		"sublime-text": {
			"wget -qO - https://download.sublimetext.com/sublimehq-pub.gpg | sudo apt-key add -",
			"sudo apt install -y apt-transport-https",
			"echo \"deb https://download.sublimetext.com/ apt/stable/\" | sudo tee /etc/apt/sources.list.d/sublime-text.list",
			"sudo apt update",
			"sudo apt install -y sublime-text",
		},
		"goland": {
			"sudo snap install goland --classic",
		},
		"webstorm": {
			"sudo snap install webstorm --classic",
		},
		"pycharm": {
			"sudo snap install pycharm-community --classic",
		},
		"intellij-idea": {
			"sudo snap install intellij-idea-community --classic",
		},
		"docker": {
			"sudo apt install -y apt-transport-https ca-certificates curl software-properties-common",
			"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -",
			"sudo add-apt-repository \"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\"",
			"sudo apt update",
			"sudo apt install -y docker-ce",
		},
	},
}

// Маппинг имен программ к пакетам для разных ОС
var packageNames = map[string]map[string]string{
	"windows": {
		"node":          "nodejs",
		"npm":           "npm",
		"yarn":          "yarn",
		"code":          "vscode",
		"webstorm":      "webstorm",
		"sublime-text":  "sublimetext3",
		"java":          "openjdk",
		"maven":         "maven",
		"gradle":        "gradle",
		"intellij-idea": "intellijidea-community",
		"eclipse":       "eclipse",
		"netbeans":      "netbeans",
		"go":            "golang",
		"python3":       "python3",
		"pip3":          "pip",
		"virtualenv":    "virtualenv",
		"git":           "git",
		"docker":        "docker-desktop",
		"curl":          "curl",
		"zsh":           "zsh",
		"jq":            "jq",
		"postman":       "postman",
		"neovim":        "neovim",
		"goland":        "goland",
		"pycharm":       "pycharm-community",
	},
	"darwin": {
		"node":          "node",
		"npm":           "npm",
		"yarn":          "yarn",
		"code":          "--cask visual-studio-code",
		"webstorm":      "--cask webstorm",
		"sublime-text":  "--cask sublime-text",
		"java":          "openjdk",
		"maven":         "maven",
		"gradle":        "gradle",
		"intellij-idea": "--cask intellij-idea-ce",
		"eclipse":       "--cask eclipse-java",
		"netbeans":      "--cask netbeans",
		"go":            "go",
		"python3":       "python3",
		"pip3":          "python3-pip",
		"virtualenv":    "virtualenv",
		"git":           "git",
		"docker":        "--cask docker",
		"curl":          "curl",
		"zsh":           "zsh",
		"jq":            "jq",
		"postman":       "--cask postman",
		"neovim":        "neovim",
		"goland":        "--cask goland",
		"pycharm":       "--cask pycharm-ce",
	},
	"linux": {
		"node":          "nodejs",
		"npm":           "npm",
		"yarn":          "yarn",
		"code":          "code",
		"webstorm":      "webstorm",
		"sublime-text":  "sublime-text",
		"java":          "openjdk-11-jdk",
		"maven":         "maven",
		"gradle":        "gradle",
		"intellij-idea": "intellij-idea-community",
		"eclipse":       "eclipse",
		"netbeans":      "netbeans",
		"go":            "golang",
		"python3":       "python3",
		"pip3":          "python3-pip",
		"virtualenv":    "python3-virtualenv",
		"git":           "git",
		"docker":        "docker.io",
		"curl":          "curl",
		"zsh":           "zsh",
		"jq":            "jq",
		"postman":       "postman",
		"neovim":        "neovim",
		"goland":        "goland",
		"pycharm":       "pycharm-community",
	},
}

func init() {
	// Настраиваем логирование
	if runtime.GOOS == "windows" {
		logFile, err := os.OpenFile("install.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(io.MultiWriter(os.Stdout, logFile))
		}
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}
}

// getPackageManager возвращает подходящий пакетный менеджер для текущей ОС
func getPackageManager(osType string) (string, error) {
	switch osType {
	case "windows":
		cmd := exec.Command("powershell", "-Command", "Get-Command choco -ErrorAction SilentlyContinue")
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("Chocolatey не установлен. Установите его выполнив следующую команду в PowerShell с правами администратора:\nSet-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))")
		}
		return "choco", nil
	case "darwin":
		if _, err := exec.Command("which", "brew").Output(); err == nil {
			return "brew", nil
		}
		return "", fmt.Errorf("Homebrew не установлен. Установите его с https://brew.sh/")
	case "linux":
		for _, pm := range []string{"apt", "yum", "dnf", "pacman"} {
			if _, err := exec.Command("which", pm).Output(); err == nil {
				return pm, nil
			}
		}
		return "", fmt.Errorf("не найден поддерживаемый пакетный менеджер")
	}
	return "", fmt.Errorf("неподдерживаемая операционная система: %s", osType)
}

// runCommand выполняет команду в системе
func runCommand(command string, osType string) error {
	log.Printf("Выполнение команды: %s\n", command)

	var cmd *exec.Cmd
	if osType == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	// Выводим команду в реальном времени
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания pipe для stdout: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("ошибка создания pipe для stderr: %v", err)
	}

	// Запускаем команду
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ошибка запуска команды: %v", err)
	}

	// Читаем вывод в реальном времени
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Printf("stdout: %s\n", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Printf("stderr: %s\n", scanner.Text())
		}
	}()

	// Ждем завершения команды
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ошибка выполнения команды: %v", err)
	}

	log.Printf("Команда выполнена успешно: %s\n", command)
	return nil
}

// detectOS определяет текущую операционную систему
func detectOS() string {
	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	default:
		return "unknown"
	}
}

// isInstalled проверяет, установлен ли инструмент
func isInstalled(program, osType string) bool {
	var cmd *exec.Cmd
	if osType == "windows" {
		// Используем PowerShell для проверки установленных программ
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Get-Command %s -ErrorAction SilentlyContinue", program))
	} else {
		cmd = exec.Command("which", program)
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Программа %s не найдена: %v\n", program, err)
		return false
	}

	log.Printf("Программа %s найдена\n", program)
	return true
}

// executeCommand выполняет команду пакетного менеджера
func executeCommand(osType, command, program string) error {
	log.Printf("Выполнение команды для ОС %s: команда=%s, программа=%s\n", osType, command, program)

	// Проверяем наличие специальных команд установки
	if command == "install" {
		if commands, exists := specialInstallCommands[osType][program]; exists {
			log.Printf("Найдены специальные команды установки для %s\n", program)
			for _, cmd := range commands {
				if err := runCommand(cmd, osType); err != nil {
					return fmt.Errorf("ошибка выполнения специальной команды для %s: %v", program, err)
				}
			}
			return nil
		}
	}

	pm, err := getPackageManager(osType)
	if err != nil {
		return err
	}

	packageName := packageNames[osType][program]
	if packageName == "" {
		packageName = program
	}

	pmCommand := packageManagerCommands[pm][command]
	if pmCommand == "" {
		return fmt.Errorf("неподдерживаемая команда %s для пакетного менеджера %s", command, pm)
	}

	var fullCommand string
	switch osType {
	case "windows":
		fullCommand = fmt.Sprintf("powershell -Command \"%s %s %s\"", pm, pmCommand, packageName)
	case "darwin":
		fullCommand = fmt.Sprintf("%s %s %s", pm, pmCommand, packageName)
	default:
		fullCommand = fmt.Sprintf("sudo apt %s %s", pmCommand, packageName)
	}

	log.Printf("Сформирована команда: %s\n", fullCommand)
	return runCommand(fullCommand, osType)
}

// installStack устанавливает все инструменты для выбранного стека
func installStack(stack Stack, ide []string, tools []string, osType string) error {
	// Проверяем права администратора для Windows
	if osType == "windows" && !checkAdminRights(osType) {
		return fmt.Errorf("необходимо запустить программу с правами администратора")
	}

	// Устанавливаем IDE
	for _, editor := range ide {
		if tool, ok := availableTools[editor]; ok {
			fmt.Printf("Установка %s...\n", tool.Description)
			if err := tool.install(osType); err != nil {
				return fmt.Errorf("ошибка установки %s: %v", editor, err)
			}
			log.Printf("IDE %s успешно установлена\n", editor)
		}
	}

	// Устанавливаем дополнительные инструменты
	for _, toolName := range tools {
		if tool, ok := availableTools[toolName]; ok {
			fmt.Printf("Установка %s...\n", tool.Description)
			if err := tool.install(osType); err != nil {
				return fmt.Errorf("ошибка установки %s: %v", toolName, err)
			}
			log.Printf("Инструмент %s успешно установлен\n", toolName)
		}
	}

	return nil
}

// checkAdminRights проверяет права администратора
func checkAdminRights(osType string) bool {
	if osType == "windows" {
		cmd := exec.Command("powershell", "-Command", "[bool](([System.Security.Principal.WindowsIdentity]::GetCurrent()).groups -match \"S-1-5-32-544\")")
		output, err := cmd.Output()
		if err != nil {
			log.Printf("Ошибка проверки прав администратора: %v\n", err)
			return false
		}
		isAdmin := strings.TrimSpace(string(output)) == "True"
		log.Printf("Права администратора: %v\n", isAdmin)
		return isAdmin
	}
	return true
}

// Вспомогательные функции для Windows
func getWindowsProgramList() []string {
	cmd := exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_Product | Select-Object Name")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Ошибка получения списка программ: %v\n", err)
		return nil
	}
	return strings.Split(string(output), "\n")
}

func isWindowsProgramInstalled(programName string) bool {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-WmiObject -Class Win32_Product | Where-Object { $_.Name -like '*%s*' }", programName))
	err := cmd.Run()
	return err == nil
}

func ensureWindowsPrerequisites() error {
	// Проверяем и включаем Windows Features, необходимые для работы
	prerequisites := []string{
		"Microsoft-Windows-Subsystem-Linux",
		"VirtualMachinePlatform",
	}

	for _, feature := range prerequisites {
		cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Enable-WindowsOptionalFeature -Online -FeatureName %s -NoRestart", feature))
		if err := cmd.Run(); err != nil {
			log.Printf("Предупреждение: не удалось включить функцию Windows %s: %v\n", feature, err)
		}
	}

	return nil
}
