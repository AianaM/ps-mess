package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"
)

type commit struct {
	Branch, Hash, Subject, Body, Author string
	AuthorDate                          string `json:"author-date"`
}

func pretty() string {
	type keyValue struct {
		key, value string
	}
	scheme := make(map[string]keyValue)
	scheme["hash"] = keyValue{key: "hash", value: "%h"}
	scheme["subject"] = keyValue{key: "subject", value: "%s"}
	scheme["body"] = keyValue{key: "body", value: "%b"}
	scheme["author"] = keyValue{key: "author", value: "%an"}
	scheme["author-date"] = keyValue{key: "author-date", value: "%ai"}
	scheme["commit-notes"] = keyValue{key: "commit-notes", value: "%N"}

	str := []string{}
	for _, v := range scheme {
		str = append(str, "[quote-here]"+v.key+"[quote-here]: [quote-here]"+v.value+"[quote-here]")
	}

	return "--pretty=format:\"{" + strings.Join(str, ",") + "},\""
}

func getLogCommands(branches []string, since, outputDir, logFile string) string {
	quoteEscape := "| sed 's/[\\\"\\'\\'']/\\\\\"/g'"
	quoteAdd := "| sed 's/\\[quote-here\\]/\"/g'"
	deleteLastComma := "| sed '$ s/.$//'"
	logs := make([]string, len(branches))
	for i, v := range branches {
		originBranch := "origin/" + v
		gitCmd := []string{"echo", "\"[\"", "$(", "git", "log", originBranch, pretty(), "--since=\"" + since + "\"", deleteLastComma, quoteEscape, quoteAdd, ")", "\"]\"", ">", outputDir + "/" + v + "." + logFile}
		logs[i] = strings.Join(gitCmd, " ")
	}
	return strings.Join(logs, " && ")
}

func getLogs(value struct{ name, path string }) []commit {
	b, err := os.ReadFile(value.path)
	if err != nil {
		log.Fatal(err)
	}
	u := []commit{}
	if err := json.NewDecoder(bytes.NewBuffer(b)).Decode(&u); err != nil {
		log.Fatal(err)
	}
	for i := range u {
		u[i].Branch = value.name
	}
	return u
}

type branch struct {
	name  string
	tasks map[string][]commit
}

type branchesLog struct {
	value map[string]branch
	tasks map[string][]commit
}

func getBranchLogs(commits []commit, branches, queueKeys []string) branchesLog {
	taskRe := searchTasksRe(queueKeys)
	b := branchesLog{
		value: make(map[string]branch, len(branches)),
		tasks: make(map[string][]commit),
	}
	for _, v := range branches {
		b.value[v] = branch{name: v, tasks: make(map[string][]commit)}
	}

	for _, v := range commits {
		found := taskRe.FindAllString(v.Subject, -1)
		for _, task := range found {
			if _, ok := b.value[v.Branch].tasks[task]; !ok {
				b.value[v.Branch].tasks[task] = []commit{v}
			} else {
				b.value[v.Branch].tasks[task] = append(b.value[v.Branch].tasks[task], v)
			}
			b.tasks[task] = append(b.tasks[task], v)
		}
	}
	return b
}

func searchTasksRe(ss []string) *regexp.Regexp {
	str := `(?im)`
	if len(ss) == 0 {
		return regexp.MustCompile(str + `((^[a-z]+)|(\s[a-z]+))-\d+`)
	}
	groups := []string{}
	for _, key := range ss {
		groups = append(groups, `(?P<`+key+`>`+strings.ToLower(key)+`-\d+)`)
	}
	return regexp.MustCompile(str + strings.Join(groups, "|"))
}

func (b branchesLog) makeCSVStr() string {
	branchLogsLen := len(b.value)
	var csv string
	csv += ";"
	branches := []string{}
	for _, v := range b.value {
		csv += v.name + ";"
		branches = append(branches, v.name)
	}
	csv += "\n"
	for task := range b.tasks {
		csv += task + ";"
		for i := 0; i < branchLogsLen; i++ {
			if _, ok := b.value[branches[i]].tasks[task]; ok {
				csv += "X;"
			} else {
				csv += ";"
			}
		}
		csv += "\n"
	}
	return csv
}
