package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputDir  = "ps-mess"
	confFile   = "config.json"
	gitLogFile = "log.json"
	tableFile  = "table.csv"
)

func init() {
	fmt.Println("Погнали!")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if dir, err := os.Getwd(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("location: ", dir)
		outputDir = dir + "/" + outputDir
	}
	if logFile, err := os.Create(outputDir + "/ps-mess.log"); err != nil {
		log.Fatalln(err)
	} else {
		log.SetOutput(logFile)
	}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	run()
}

func save(path, content string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.WriteString(content + "\n")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
	f.Sync()
}

func run() {
	if len(os.Args) < 2 {
		fmt.Println("нужно уточнение:")
		fmt.Println("argument prep - создаем файл настроек и команду для сбора логов")
		fmt.Println("argument comp - читаем файлы логов и создаем файлик с табличкой сравнения")
		log.Println("no arg")
		return
	}
	arg := os.Args[1]
	switch arg {
	case "prep":
		prep()
	case "comp":
		comp()
	default:
		fmt.Println("Не понял что надо делать, ожидается prep или comp")
		log.Println("Unexpected argument")
	}
}

func prep() {
	confPath := outputDir + "/" + confFile
	fmt.Println("Создаю файл настроек", confPath, "его можно и нужно будет изменить")
	log.Println("prep run", confPath)
	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		save(confPath, makeFakeConfig())
	}
	conf := getSettings(confPath)
	commands := getLogCommands(conf.Branches, conf.Since, outputDir, gitLogFile)
	fmt.Println("===\n\n", "Эту команду надо ввести в той дирректории где у тебя репа:", "\n", commands, "\n\n===")
	log.Println("prep done")
}

func comp() {
	log.Println("comp run")
	logs := getLogFiles(outputDir)
	commits := []commit{}
	branches := make([]string, len(logs))
	for i, v := range logs {
		commits = append(commits, getLogs(v)...)
		branches[i] = v.name
	}
	conf := getSettings(outputDir + "/" + confFile)
	bLog := getBranchLogs(commits, branches, conf.QueueKeys)
	csv := bLog.makeCSVStr()
	tablePath := outputDir + "/" + tableFile
	save(tablePath, csv)
	fmt.Println("Done! Табличка тут", tablePath)
	log.Println("comp done", tablePath)
}

// Ищем файлы *.log.json
func getLogFiles(dir string) []struct{ name, path string } {
	log.Println("getLogFiles run")
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}
	logs := []struct{ name, path string }{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		name := d.Name()
		if !d.IsDir() && strings.HasSuffix(path, ".log.json") {
			logs = append(logs, struct{ name, path string }{name: name, path: path})
			return nil
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("getLogFiles done, found ", len(logs), "files")
	return logs
}
