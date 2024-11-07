package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type brunchLog struct {
	name, path                                                string
	tasks                                                     []string
	tasksCommits, intersection, leftExclusive, rightExclusive map[string][]string
}

func newBrunchLog(name string) brunchLog {
	return brunchLog{name: name, path: "src/log-" + name + ".txt"}
}

func main() {
	start()
}

func start() {
	logs := []string{"dev", "test", "main"}
	cols := []brunchLog{}
	for _, v := range logs {
		l := newBrunchLog(v)
		l.open()
		cols = append(cols, l)
	}
	table := comparisonTable(cols)
	save("src/comparing-table.csv", table)
	tableSimple := comparisonTableSimple(cols)
	save("src/comparing-table-simple.csv", tableSimple)
}

func (l *brunchLog) open() {
	file, err := os.Open(l.path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	l.tasksCommits = make(map[string][]string)

	r := bufio.NewReader(file)
	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			str := string(line)
			tasks := searchPStasks(str)
			for _, v := range tasks {
				l.tasksCommits[v] = append(l.tasksCommits[v], str)
			}
		}
		if err != nil {
			break
		}
	}
	l.tasks = keys(l.tasksCommits)
}

func (l *brunchLog) compare(a *brunchLog) {
	left := l.tasks
	right := a.tasks
	intersection := []string{}
	for i, l := range left {
		for ii, r := range right {
			if l == r {
				intersection = append(intersection, r)
				left = append(left[:i], left[i+1:]...)
				right = append(right[:ii], right[ii+1:]...)
			}
		}
	}
	l.intersection[a.path] = intersection
	l.leftExclusive[a.path] = left

	a.intersection[l.path] = intersection
	a.leftExclusive[l.path] = right
}
func comparisonTableSimple(logs []brunchLog) string {
	filled := func(n int) []string {
		s := make([]string, n)
		for i := range s {
			s[i] = "-"
		}
		return s
	}
	colsLen := len(logs)
	cols := filled(colsLen)
	rows := []string{}

	for lCol, lLog := range logs {
		cols[lCol] = lLog.path
		for _, lTask := range lLog.tasks {
			row := filled(colsLen)
			row[lCol] = lTask
			for rCol := lCol + 1; rCol < colsLen; rCol++ {
				rLog := logs[rCol]
				for ii, rTask := range rLog.tasks {
					if lTask == rTask {
						row[rCol] = "+"
						logs[rCol].tasks = append(logs[rCol].tasks[:ii], logs[rCol].tasks[ii+1:]...)
						break
					}

				}
			}
			rows = append(rows, strings.Join(row, ";"))
		}
	}
	return strings.Join(cols, ";") + "\n" + strings.Join(rows, "\n")
}

func comparisonTable(logs []brunchLog) string {
	colsLen := len(logs)
	cols := make([]string, colsLen)
	rows := []string{}

	for lCol, lLog := range logs {
		cols[lCol] = lLog.path
		for _, lTask := range lLog.tasks {
			row := make([]string, colsLen)
			row[lCol] = lTask + "( " + strings.Join(lLog.tasksCommits[lTask], "----> ") + " )"
			for rCol := lCol + 1; rCol < colsLen; rCol++ {
				rLog := logs[rCol]
				for ii, rTask := range rLog.tasks {
					if lTask == rTask {
						row[rCol] = strings.Join(logs[rCol].tasksCommits[rTask], "----> ")
						logs[rCol].tasks = append(logs[rCol].tasks[:ii], logs[rCol].tasks[ii+1:]...)
						break
					}

				}
			}
			rows = append(rows, strings.Join(row, ";"))
		}
	}
	return strings.Join(cols, ";") + "\n" + strings.Join(rows, "\n")
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

func searchPStasks(str string) []string {
	re := regexp.MustCompile(`(?im)(?P<ps>ps-\d+)`)
	return re.FindAllString(str, -1)
}

func keys(m map[string][]string) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}
