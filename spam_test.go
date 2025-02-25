package main

import (
	"testing"
)

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

func Test_hasHighEntropy(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Spam - l6bh77d756344w1v2222n8m",
			args: args{"l6bh77d756344w1v2222n8m"},
			want: true,
		},
		{
			name: "Spam - 6zml0l5dk55m5l3767tcj5o",
			args: args{"6zml0l5dk55m5l3767tcj5o"},
			want: true,
		},
		{
			name: "Spam - c3m1e9mav536u973wc16bqu",
			args: args{"c3m1e9mav536u973wc16bqu"},
			want: true,
		},
		{
			name: "Spam - nwyel400g4in5l115nu21jf",
			args: args{"nwyel400g4in5l115nu21jf"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasHighEntropy(tt.args.s); got != tt.want {
				t.Errorf("hasHighEntropy() = %v, want %v", got, tt.want)
			}
		})
	}
}
