package helper_test

import (
	"testing"

	"github.com/i5heu/bonito-cache/internal/helper"
)

func TestSanitizeMimeType(t *testing.T) {
	type args struct {
		mime string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty string",
			args: args{
				mime: "",
			},
			want: "application/octet-stream",
		},
		{
			name: "long string",
			args: args{
				mime: "loremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsum/loremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsumloremipsum",
			},
			want: "application/octet-stream",
		},
		{
			name: "non ascii",
			args: args{
				mime: "video/mp4💖",
			},
			want: "application/octet-stream",
		},
		{
			name: "normal video",
			args: args{
				mime: "video/mp4",
			},
			want: "video/mp4",
		},
		{
			name: "normal img",
			args: args{
				mime: "image/jpeg",
			},
			want: "image/jpeg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := helper.SanitizeMimeType(tt.args.mime); got != tt.want {
				t.Errorf("SanitizeMimeType() = %v, want %v", got, tt.want)
			}
		})
	}
}
