package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type config struct {
	// Ключи очередей в яндекс треккере
	QueueKeys []string
	// Ветки, которые нас интересуют
	Branches []string
	// За какой период собираем логи
	Since string
}

func makeFakeConfig() string {
	log.Println("makeFakeConfig start")
	consf := config{
		QueueKeys: []string{"ps", "scp"},
		Branches:  []string{"dev", "test"},
		Since:     "3 weeks ago",
	}
	b := new(strings.Builder)
	encoder := json.NewEncoder(b)
	if err := encoder.Encode(consf); err != nil {
		fmt.Println("ошибка при создании файла примера конфигов")
		log.Fatal(err)
	}
	log.Println("makeFakeConfig done")
	return b.String()
}

func getSettings(path string) config {
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("не могу открыть файл ", path)
		log.Fatal(err)
	}
	var u config
	json.NewDecoder(bytes.NewBuffer(b)).Decode(&u)
	return u
}
