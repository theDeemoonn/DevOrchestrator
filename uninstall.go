// uninstall.go
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Удалить инструменты",
	Run:   uninstallRun,
}

func uninstallRun(cmd *cobra.Command, args []string) {
	osType := detectOS()
	fmt.Printf("Обнаруженная ОС: %s\n", osType)

	stackStr := selectStack()
	tools := selectStackTools(stackStr)

	var wg sync.WaitGroup
	errorsCh := make(chan error, len(tools))

	for _, toolName := range tools {
		if tool, ok := availableTools[toolName]; ok {
			wg.Add(1)
			go func(tool Tool) {
				defer wg.Done()
				err := tool.uninstall(osType)
				if err != nil {
					errorsCh <- fmt.Errorf("ошибка удаления %s: %v", tool.Description, err)
				}
			}(tool)
		}
	}

	wg.Wait()
	close(errorsCh)

	for err := range errorsCh {
		log.Printf("Ошибка: %v\n", err)
	}
}
