package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Переменные для команды git log
type gitLog struct {
	// Путь до папки с гитом
	Path string
	// Ветки, которые в которые сравнивать
	Branches []string
	// С каких пор собираем логи
	Since string
}

func (c gitLog) save(outputDir string) {
	log.Default().Println("Пошел собирать логи")
	// Текущая рабочая директория
	root, _ := os.Getwd()
	fails := []string{}
	for _, v := range c.Branches {
		filename := "log-" + v + ".txt"
		output := []string{root, outputDir, filename}
		cmd := exec.Command("git", "log", "origin/"+v, "--pretty=format:\"%h%x09%s %an%x09%ad%x09%n\"", "--since=\""+c.Since+"\"", "--output=\""+strings.Join(output, "/")+"\"")
		cmd.Dir = c.Path

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			log.Println("git log err:", err)
			fails = append(fails, cmd.String())
		}
		fmt.Println(out.String())
	}
	fmt.Println("Собрал логи, по крайней мере попытался")

	if len(fails) > 0 {
		fmt.Println("...")
		fmt.Println("Не смог получить все логи, придется самостоятельно положить логи в папку", root+outputDir)
		fmt.Println(strings.Join(fails, " && "))
		fmt.Println("...")

		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			log.Fatal(err)
		}
	}
}
