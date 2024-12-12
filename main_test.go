package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_convertIDtoCreationTime(t *testing.T) {
	timestamp := "2017-01-28 06:02:39.924 +0000 UTC"
	timeObj, _ := time.Parse("2006-01-02 15:04:05.999 -0700 MST", timestamp)

	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "274781269861072896",
			args: args{"274781269861072896"},
			want: timeObj,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertIDtoCreationTime(tt.args.id).UTC() // Convert to UTC
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertIDtoCreationTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUsernameSpam(t *testing.T) {
	patterns := []string{
		`(?i)Announcement`,
		`(?i)FAQ`,
		`\d{4}$`,
	}

	type args struct {
		username string
		patterns []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Ham - janedoe",
			args: args{"janedoe", patterns},
			want: false,
		},
		{
			name: "Spam - announcements23",
			args: args{"announcements23", patterns},
			want: true,
		},
		{
			name: "Spam - spammer1234",
			args: args{"spammer1234", patterns},
			want: true,
		},
		{
			name: "Ham - user123",
			args: args{"user123", patterns},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNameSpam(tt.args.username, tt.args.patterns); got != tt.want {
				t.Errorf("isUsernameSpam() = %v, want %v", got, tt.want)
			}
		})
	}
}
