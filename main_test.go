package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
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
			re := searchTasksRe(tt.conf.QueueKeys)
			got := re.FindAllString(tt.str, -1)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchTasksRe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLogFiles(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		files   fstest.MapFS
		want    []struct{ name, path string }
		wantErr bool
	}{
		{
			name: "Single log file",
			dir:  "testdir",
			files: fstest.MapFS{
				"testdir/file1.log.json": &fstest.MapFile{Mode: fs.ModePerm},
			},
			want: []struct{ name, path string }{
				{name: "file1.log.json", path: "testdir/file1.log.json"},
			},
			wantErr: false,
		},
		{
			name: "Multiple log files",
			dir:  "testdir",
			files: fstest.MapFS{
				"testdir/file1.log.json": &fstest.MapFile{Mode: fs.ModePerm},
				"testdir/file2.log.json": &fstest.MapFile{Mode: fs.ModePerm},
			},
			want: []struct{ name, path string }{
				{name: "file1.log.json", path: "testdir/file1.log.json"},
				{name: "file2.log.json", path: "testdir/file2.log.json"},
			},
			wantErr: false,
		},
		{
			name: "No log files",
			dir:  "testdir",
			files: fstest.MapFS{
				"testdir/file1.txt": &fstest.MapFile{Mode: fs.ModePerm},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Mixed files",
			dir:  "testdir",
			files: fstest.MapFS{
				"testdir/file1.log.json": &fstest.MapFile{Mode: fs.ModePerm},
				"testdir/file2.txt":      &fstest.MapFile{Mode: fs.ModePerm},
			},
			want: []struct{ name, path string }{
				{name: "file1.log.json", path: "testdir/file1.log.json"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for the test
			tmpDir := t.TempDir()
			for path, file := range tt.files {
				fullPath := filepath.Join(tmpDir, path)
				if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte{}, file.Mode); err != nil {
					t.Fatalf("Failed to create file: %v", err)
				}
			}

			got := getLogFiles(tmpDir)
			var gotFiles []struct{ name, path string }
			for _, file := range got {
				path := strings.TrimPrefix(file.path, tmpDir+"/")
				gotFiles = append(gotFiles, struct{ name, path string }{name: file.name, path: path})
			}
			if !reflect.DeepEqual(gotFiles, tt.want) {
				t.Errorf("getLogFiles() = %v, want %v", gotFiles, tt.want)
			}
		})
	}
}
