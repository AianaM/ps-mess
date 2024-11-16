package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type branchLog struct {
	name, path string
	tasks      map[string][]string
}

func newBranchLog(name, path string, tasks map[string][]string) branchLog {
	return branchLog{name: name, path: path, tasks: tasks}
}

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func init() {
	setLogger()
	logger.Info("Погнали!")
}
func main() {
	confPath := "src/config.json"
	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		confPath = configPath()
	}
	conf := getSettings(confPath)
	if conf.Auto {
		conf.Cmd.save(conf.Src)
	}
	branchLogs, tasks := getBranchLogs(conf)
	save(conf.Output+"/comparing-table-simple.csv", toCSVstr(branchLogs, tasks))
}

func toCSVstr(branchLogs []branchLog, tasks map[string]bool) string {
	branchLogsLen := len(branchLogs)
	var csv string
	csv += ";"
	for i := 0; i < branchLogsLen; i++ {
		csv += branchLogs[i].name + ";"
	}
	csv += "\n"
	for task := range tasks {
		csv += task + ";"
		for i := 0; i < branchLogsLen; i++ {
			if _, ok := branchLogs[i].tasks[task]; ok {
				csv += "X;"
			} else {
				csv += ";"
			}
		}
		csv += "\n"
	}
	return csv
}

func getBranchLogs(conf config) (branchLogs []branchLog, tasks map[string]bool) {
	logs := getLogs(conf.Src)
	re := searchTasksRe(conf)
	var wg sync.WaitGroup

	tasks = make(map[string]bool)
	tasksMx := sync.RWMutex{}
	search := func(branch branchLog) func(str string) {
		return func(str string) {
			res := re.FindAllString(str, -1)
			for _, v := range res {
				tasksMx.RLock()
				_, ok := tasks[v]
				tasksMx.RUnlock()
				if !ok {
					tasksMx.Lock()
					tasks[v] = true
					tasksMx.Unlock()
				}
				branch.tasks[v] = append(branch.tasks[v], str)
			}
		}
	}
	branchLogs = make([]branchLog, len(logs))
	for i, v := range logs {
		wg.Add(1)
		b := newBranchLog(v.name, v.path, make(map[string][]string))
		branchLogs[i] = b
		go getTasks(v.path, search(b), &wg)
	}

	wg.Wait()

	return branchLogs, tasks
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
