package main

import (
	"reflect"
	"testing"
)

func TestGetBranchLogs(t *testing.T) {
	commits := []commit{
		{Branch: "branch1", Subject: "id 0 task-123: Fix bug", Body: "Fixed a bug", Author: "Author1", AuthorDate: "2023-01-01"},
		{Branch: "branch2", Subject: "id 1 task-456: Add feature", Body: "Added a new feature", Author: "Author2", AuthorDate: "2023-01-02"},
		{Branch: "branch1", Subject: "id 2 task-123: Improve performance", Body: "Improved performance", Author: "Author1", AuthorDate: "2023-01-03"},
		{Branch: "branch3", Subject: "id 3 task-789 & task-790/task-791,task-123: Update documentation", Body: "Updated the documentation", Author: "Author3", AuthorDate: "2023-01-04"},
		{Branch: "branch2", Subject: "id 4 task-456: Refactor code", Body: "Refactored the code", Author: "Author2", AuthorDate: "2023-01-05"},
		{Branch: "branch3", Subject: "id 5 task-791: Fix typo", Body: "Fixed a typo", Author: "Author3", AuthorDate: "2023-01-06"},
	}
	branches := []string{"branch1", "branch2", "branch3"}
	queueKeys := []string{"task"}

	expected := branchesLog{
		value: map[string]branch{
			"branch1": {
				name: "branch1",
				tasks: map[string][]commit{
					"task-123": {commits[0], commits[2]},
				},
			},
			"branch2": {
				name: "branch2",
				tasks: map[string][]commit{
					"task-456": {commits[1], commits[4]},
				},
			},
			"branch3": {
				name: "branch3",
				tasks: map[string][]commit{
					"task-789": {commits[3]},
					"task-790": {commits[3]},
					"task-791": {commits[3], commits[5]},
					"task-123": {commits[3]},
				},
			},
		},
		tasks: map[string][]commit{
			"task-123": {commits[0], commits[2], commits[3]},
			"task-456": {commits[1], commits[4]},
			"task-789": {commits[3]},
			"task-790": {commits[3]},
			"task-791": {commits[3], commits[5]},
		},
	}

	result := getBranchLogs(commits, branches, queueKeys)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
