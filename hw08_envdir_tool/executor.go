package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...)

	// Устанавливаем стандартный вывод и стандартный вывод ошибок
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	commandEnv := []string{}
	for envItemName, envItemValue := range env {
		commandEnv = append(commandEnv, envItemName+"="+envItemValue.Value)
	}
	command.Env = commandEnv
	// Запускаем команду
	if err := command.Run(); err != nil {
		fmt.Printf("Ошибка при запуске приложения: %v\n", err)
	}
	return command.ProcessState.ExitCode()
}
