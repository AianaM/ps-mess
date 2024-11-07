package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type brunchLog struct {
	path                                        string
	tasks                                       []string
	intersection, leftExclusive, rightExclusive map[string][]string
}

func newBrunchLog(path string) brunchLog {
	return brunchLog{path: path, tasks: []string{}}
}

func main() {
	start()
}

func start() {
	logs := []string{"src/log-dev.txt", "src/log-main.txt", "src/log-test.txt"}
	cols := []brunchLog{}
	for _, v := range logs {
		l := newBrunchLog(v)
		l.open()
		cols = append(cols, l)
	}
	lines := comparisonTable(cols)
	for _, v := range lines {
		fmt.Println(v)
	}
}

func (l *brunchLog) open() {
	file, err := os.Open(l.path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			l.tasks = append(l.tasks, searchPStasks(string(line))...)
		}
		if err != nil {
			break
		}
	}
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

func comparisonTable(logs []brunchLog) [][]string {
	filled := func(n int) []string {
		s := make([]string, n)
		for i := range s {
			s[i] = "-"
		}
		return s
	}

	colsLen := len(logs)
	lines := [][]string{filled(colsLen)}

	for lCol, lLog := range logs {
		lines[0][lCol] = lLog.path
		for _, lTask := range lLog.tasks {
			line := filled(colsLen)
			line[lCol] = lTask
			for rCol := lCol + 1; rCol < colsLen; rCol++ {
				rLog := logs[rCol]
				for ii, rTask := range rLog.tasks {
					if lTask == rTask {
						line[rCol] = "+"
						logs[rCol].tasks = append(logs[rCol].tasks[:ii], logs[rCol].tasks[ii+1:]...)
						break
					}

				}
			}
			lines = append(lines, line)
		}
	}
	return lines
}

func searchPStasks(str string) []string {
	re := regexp.MustCompile(`(?im)(?P<ps>ps-\d+)`)
	return re.FindAllString(str, -1)
}
