package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Обновить инструменты",
	Run:   updateRun,
}

func updateRun(cmd *cobra.Command, args []string) {
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
				err := tool.update(osType)
				if err != nil {
					errorsCh <- fmt.Errorf("ошибка обновления %s: %v", tool.Description, err)
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
