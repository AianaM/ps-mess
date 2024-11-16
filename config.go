package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type config struct {
	// Путь к файлам логов
	Src string
	// Выходной каталог.
	Output string
	// Ключи очередей в яндекс треккере
	QueueKeys []string
	Cmd       gitLog
	// Надо ли собирать логи (если false, то надо самостоятельно положить логи в папку указанную в src)
	Auto bool
}

func configPath() string {
	fmt.Println("Путь до файла с конфигами, например: src/config.json")
	fmt.Print("\n")
	fmt.Print("Вводи: ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}
	return input
}

func getSettings(path string) config {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read file: %v\n", err)
	}
	var u config
	json.NewDecoder(bytes.NewBuffer(b)).Decode(&u)
	return u
}
