package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSearchTasksRe(t *testing.T) {
	tests := []struct {
		name string
		conf config
		str  string
		want []string
	}{
		{
			name: "Single queue key",
			conf: config{QueueKeys: []string{"ps"}},
			str:  "This string contains ps-123 task.",
			want: []string{"ps-123"},
		},
		{
			name: "Multiple queue keys",
			conf: config{QueueKeys: []string{"ps", "ts"}},
			str:  "This string contains ps-123 and ts-456 tasks.",
			want: []string{"ps-123", "ts-456"},
		},
		{
			name: "No matching tasks",
			conf: config{QueueKeys: []string{"ps"}},
			str:  "This string contains no tasks.",
			want: nil,
		},
		{
			name: "Tasks with different cases",
			conf: config{QueueKeys: []string{"ps"}},
			str:  "This string contains PS-123 and ps-456 tasks.",
			want: []string{"PS-123", "ps-456"},
		},
		{
			name: "Mixed content",
			conf: config{QueueKeys: []string{"ps", "ts"}},
			str:  "ps-123, some text, ts-456, more text, ps-789.",
			want: []string{"ps-123", "ts-456", "ps-789"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			re := searchTasksRe(tt.conf)
			got := re.FindAllString(tt.str, -1)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchTasksRe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogs(t *testing.T) {
	testCases := []struct {
		name     string
		logFiles []string
		expected map[string]bool
	}{
		{
			name:     "Three log files",
			logFiles: []string{"log-1.txt", "log-2.txt", "log-3.txt"},
			expected: map[string]bool{
				"log-1.txt": false,
				"log-2.txt": false,
				"log-3.txt": false,
			},
		},
		{
			name:     "No log files",
			logFiles: []string{},
			expected: map[string]bool{},
		},
		{
			name:     "Mixed files",
			logFiles: []string{"log-1.txt", "log-2.txt", "otherfile.txt"},
			expected: map[string]bool{
				"log-1.txt": false,
				"log-2.txt": false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir, err := ioutil.TempDir("", "testlogs")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			// Create test log files
			for _, fileName := range tc.logFiles {
				filePath := filepath.Join(tempDir, fileName)
				if err := ioutil.WriteFile(filePath, []byte("test content"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			// Call getLogs function
			logs := getLogs(tempDir)

			// Verify the returned logs
			for _, log := range logs {
				if _, ok := tc.expected[log.name]; ok {
					tc.expected[log.name] = true
				} else {
					t.Errorf("Unexpected log file: %s", log.name)
				}
			}

			for fileName, found := range tc.expected {
				if !found {
					t.Errorf("Expected log file not found: %s", fileName)
				}
			}
		})
	}
}
