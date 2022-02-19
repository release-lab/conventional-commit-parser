package conventionalcommitparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_paseFooterParagraph(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name string
		args args
		want Footer
	}{
		{
			name: "invalid footer paragraph",
			args: args{
				txt: "hello world",
			},
			want: Footer{Title: "hello world"},
		},
		{
			name: "invalid footer paragraph",
			args: args{
				txt: "helloworld",
			},
			want: Footer{Title: "helloworld"},
		},
		{
			name: "tag footer paragraph",
			args: args{
				txt: "Refs: #2",
			},
			want: Footer{Tag: "Refs", Title: "#2"},
		},
		{
			name: "tag with - footer paragraph",
			args: args{
				txt: "Refs-With: #2",
			},
			want: Footer{Tag: "Refs-With", Title: "#2"},
		},
		{
			name: "tag with multiple - footer paragraph",
			args: args{
				txt: "Refs-With-User: #2",
			},
			want: Footer{Tag: "Refs-With-User", Title: "#2"},
		},
		{
			name: "hash footer paragraph",
			args: args{
				txt: "Close #1, #2",
			},
			want: Footer{Tag: "Close", Title: "#1, #2"},
		},
		{
			name: "hash footer paragraph with multiple spaces",
			args: args{
				txt: "Close   #1, #2  ",
			},
			want: Footer{Tag: "Close", Title: "#1, #2"},
		},
		{
			name: "invalid hash footer paragraph with # prefix",
			args: args{
				txt: "#Close #1, #2",
			},
			want: Footer{Tag: "", Title: "#Close #1, #2"},
		},
		{
			name: "invalid hash footer paragraph without # prefix",
			args: args{
				txt: "Close 1, #2",
			},
			want: Footer{Tag: "", Title: "Close 1, #2"},
		},
		{
			name: "BREAKING CHANGE footer paragraph",
			args: args{
				txt: "BREAKING CHANGE: this is a breaking change",
			},
			want: Footer{Tag: "BREAKING CHANGE", Title: "this is a breaking change"},
		},
		{
			name: "invalid BREAKING CHANGES footer paragraph",
			args: args{
				txt: "BREAKING CHANGES: this is a breaking change",
			},
			want: Footer{Tag: "", Title: "BREAKING CHANGES: this is a breaking change"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, paseFooterParagraph(tt.args.txt))
		})
	}
}

func Test_splitToLines(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single line",
			args: args{
				text: "single line",
			},
			want: []string{"single line"},
		},
		{
			name: "multiple line",
			args: args{
				text: "line1\nline2",
			},
			want: []string{"line1", "line2"},
		},
		{
			name: "multiple line2",
			args: args{
				text: "line1\r\nline2",
			},
			want: []string{"line1", "line2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, splitToLines(tt.args.text))
		})
	}
}

func TestParseFooter(t *testing.T) {
	type args struct {
		txt string
	}
	tests := []struct {
		name string
		args args
		want Footer
	}{
		{
			name: "invalid footer",
			args: args{
				txt: "invalid footer",
			},
			want: Footer{
				Tag:     "",
				Title:   "invalid footer",
				Content: "",
			},
		},
		{
			name: "tag: footer",
			args: args{
				txt: "tag: footer",
			},
			want: Footer{
				Tag:     "tag",
				Title:   "footer",
				Content: "",
			},
		},
		{
			name: "Closes #1",
			args: args{
				txt: "Closes #1",
			},
			want: Footer{
				Tag:     "Closes",
				Title:   "#1",
				Content: "",
			},
		},
		{
			name: "Closes #1, #2, #3",
			args: args{
				txt: "Closes #1, #2, #3",
			},
			want: Footer{
				Tag:     "Closes",
				Title:   "#1, #2, #3",
				Content: "",
			},
		},
		{
			name: "Closes   #1, #2, #3  ",
			args: args{
				txt: "Closes   #1, #2, #3  ",
			},
			want: Footer{
				Tag:     "Closes",
				Title:   "#1, #2, #3",
				Content: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseFooter(tt.args.txt))
		})
	}
}
