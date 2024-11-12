package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type brunchLog struct {
	name, path string
	tasks      map[string][]string
}

func newBrunchLog(name, path string, tasks map[string][]string) brunchLog {
	return brunchLog{name: name, path: path, tasks: tasks}
}

type config struct {
	Logs      string
	Output    string
	QueueKeys []string
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)
}

func main() {
	confPath := "src/settings.json" //configPath()
	conf := getSettings(confPath)
	brunchLogs, tasks := getBrunchLogs(conf)
	brunchLogsLen := len(brunchLogs)

	var csv string
	csv += ";"
	for i := 0; i < brunchLogsLen; i++ {
		csv += brunchLogs[i].name + ";"
	}
	csv += "\n"
	for task := range tasks {
		csv += task + ";"
		for i := 0; i < brunchLogsLen; i++ {
			if _, ok := brunchLogs[i].tasks[task]; ok {
				csv += "X;"
			} else {
				csv += ";"
			}
		}
		csv += "\n"
	}
	save(conf.Output+"/comparing-table-simple.csv", csv)
}

func getBrunchLogs(conf config) (brunchLogs []brunchLog, tasks map[string]bool) {
	logs := getLogs(conf.Logs)
	re := searchTasksRe(conf)
	var wg sync.WaitGroup

	tasks = make(map[string]bool)
	search := func(brunch brunchLog) func(str string) {
		return func(str string) {
			res := re.FindAllString(str, -1)
			for _, v := range res {
				if _, ok := tasks[v]; !ok {
					tasks[v] = true
				}
				brunch.tasks[v] = append(brunch.tasks[v], str)
			}
		}
	}
	brunchLogs = make([]brunchLog, len(logs))
	for i, v := range logs {
		wg.Add(1)
		b := newBrunchLog(v.name, v.path, make(map[string][]string))
		brunchLogs[i] = b
		go getTasks(v.path, search(b), &wg)
	}

	wg.Wait()

	return brunchLogs, tasks
}

func getTasks(path string, fn func(str string), wg *sync.WaitGroup) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		file.Close()
		wg.Done()
	}()

	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			fn(string(line))
		}
		if err != nil {
			break
		}
	}
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

func searchTasksRe(conf config) *regexp.Regexp {
	str := `(?im)`
	groups := []string{}
	for _, key := range conf.QueueKeys {
		groups = append(groups, `(?P<`+key+`>`+strings.ToLower(key)+`-\d+)`)
	}
	return regexp.MustCompile(str + strings.Join(groups, "|"))
}

func configPath() string {
	fmt.Println("Путь до файла с конфигами, например: src/settings.json")
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

// Ищем файлы с префиксом log- и расширением .txt
func getLogs(dir string) []struct{ name, path string } {
	logs := []struct{ name, path string }{}
	fmt.Println(dir)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasPrefix(name, "log-") && strings.HasSuffix(path, ".txt") {
			logs = append(logs, struct{ name, path string }{name: name, path: path})
			return nil
		} else {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return logs
}
