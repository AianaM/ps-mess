package main

import (
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
