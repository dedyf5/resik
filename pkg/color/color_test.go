// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package color

import (
	"fmt"
	"testing"
)

func Test_color(t *testing.T) {
	tests := []struct {
		name string
		c    int
		want string
	}{
		{
			name: "none color reset code",
			c:    NONE,
			want: fmt.Sprintf("%s[%dm", escape, NONE),
		},
		{
			name: "red color code",
			c:    RED,
			want: fmt.Sprintf("%s[3%dm", escape, RED),
		},
		{
			name: "green color code",
			c:    GREEN,
			want: fmt.Sprintf("%s[3%dm", escape, GREEN),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := color(tt.c); got != tt.want {
				t.Errorf("color() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		c    int
		text string
		want string
	}{
		{
			name: "format text with red color",
			c:    RED,
			text: "hello",
			want: color(RED) + "hello" + color(NONE),
		},
		{
			name: "format text with none color",
			c:    NONE,
			text: "world",
			want: color(NONE) + "world" + color(NONE),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Format(tt.c, tt.text); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
