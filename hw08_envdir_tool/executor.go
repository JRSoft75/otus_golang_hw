package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	// Устанавливаем стандартный ввод, вывод и стандартный вывод ошибок
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	// Получаем список переменных окружения
	envVars := os.Environ()

	// Создаем мапу для хранения переменных окружения
	envMap := make(map[string]string)

	// Заполняем мапу переменными окружения
	for _, envVar := range envVars {
		// Разделяем строку на ключ и значение
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			envMap[key] = value
		}
	}

	for key, newEnvItem := range env {
		if newEnvItem.NeedRemove {
			delete(envMap, key)
		} else {
			envMap[key] = newEnvItem.Value
		}
	}
	commandEnv := []string{}
	for envItemName, envItemValue := range envMap {
		commandEnv = append(commandEnv, envItemName+"="+envItemValue)
	}

	command.Env = commandEnv
	// Запускаем команду
	if err := command.Run(); err != nil {
		fmt.Printf("Ошибка при запуске приложения: %v\n", err)
	}
	return command.ProcessState.ExitCode()
}
