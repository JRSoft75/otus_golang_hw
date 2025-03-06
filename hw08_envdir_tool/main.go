package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Args содержит аргументы командной строки
	// Первый элемент - это имя исполняемого файла
	args := os.Args
	// Проверяем, есть ли аргументы
	if len(args) < 2 {
		fmt.Println("Пожалуйста, укажите аргументы командной строки.")
		return
	}

	dirPath := args[1]
	commandWithArgs := args[2:]

	env, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	os.Exit(RunCmd(commandWithArgs, env))
}
