package conventionalcommitparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		args   args
		want   Message
		header Header
		footer []Footer
	}{
		{
			name: "common commit",
			args: args{
				message: "this is a commit message",
			},
			want: Message{
				Header: "this is a commit message",
				Body:   "",
				Footer: []string{},
			},
			header: Header{
				Type:      "",
				Scope:     "",
				Subject:   "this is a commit message",
				Important: false,
			},
			footer: []Footer{},
		},
		{
			name: "common commit with body",
			args: args{
				message: "this is a commit message\n\nthis is commit body",
			},
			want: Message{
				Header: "this is a commit message",
				Body:   "this is commit body",
				Footer: []string{},
			},
			header: Header{
				Type:      "",
				Scope:     "",
				Subject:   "this is a commit message",
				Important: false,
			},
			footer: []Footer{},
		},
		{
			name: "Commit message with description and breaking change footer",
			args: args{
				message: `feat: allow provided config object to extend other configs

BREAKING CHANGE: 'extends' key in config file is now used for extending other config files`,
			},
			want: Message{
				Header: "feat: allow provided config object to extend other configs",
				Body:   "",
				Footer: []string{"BREAKING CHANGE: 'extends' key in config file is now used for extending other config files"},
			},
			header: Header{
				Type:      "feat",
				Scope:     "",
				Subject:   "allow provided config object to extend other configs",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:     "BREAKING CHANGE",
					Title:   "'extends' key in config file is now used for extending other config files",
					Content: "",
				},
			},
		},
		{
			name: "refactor!: drop support for Node 6",
			args: args{
				message: `refactor!: drop support for Node 6`,
			},
			want: Message{
				Header: "refactor!: drop support for Node 6",
				Body:   "",
				Footer: []string{},
			},
			header: Header{
				Type:      "refactor",
				Scope:     "",
				Subject:   "drop support for Node 6",
				Important: true,
			},
			footer: []Footer{},
		},
		{
			name: "Commit message with scope and ! to draw attention to breaking change",
			args: args{
				message: `refactor(runtime)!: drop support for Node 6`,
			},
			want: Message{
				Header: "refactor(runtime)!: drop support for Node 6",
				Body:   "",
				Footer: []string{},
			},
			header: Header{
				Type:      "refactor",
				Scope:     "runtime",
				Subject:   "drop support for Node 6",
				Important: true,
			},
			footer: []Footer{},
		},
		{
			name: "Commit message with both ! and BREAKING CHANGE footer",
			args: args{
				message: `refactor!: drop support for Node 6

BREAKING CHANGE: refactor to use JavaScript features not available in Node 6.`,
			},
			want: Message{
				Header: "refactor!: drop support for Node 6",
				Body:   "",
				Footer: []string{"BREAKING CHANGE: refactor to use JavaScript features not available in Node 6."},
			},
			header: Header{
				Type:      "refactor",
				Scope:     "",
				Subject:   "drop support for Node 6",
				Important: true,
			},
			footer: []Footer{
				{
					Tag:     "BREAKING CHANGE",
					Title:   "refactor to use JavaScript features not available in Node 6.",
					Content: "",
				},
			},
		},
		{
			name: "Commit message with no body",
			args: args{
				message: `docs: correct spelling of CHANGELOG`,
			},
			want: Message{
				Header: "docs: correct spelling of CHANGELOG",
				Body:   "",
				Footer: []string{},
			},
			header: Header{
				Type:      "docs",
				Scope:     "",
				Subject:   "correct spelling of CHANGELOG",
				Important: false,
			},
			footer: []Footer{},
		},
		{
			name: "Commit message with scope",
			args: args{
				message: `feat(lang): add polish language`,
			},
			want: Message{
				Header: "feat(lang): add polish language",
				Body:   "",
				Footer: []string{},
			},
			header: Header{
				Type:      "feat",
				Scope:     "lang",
				Subject:   "add polish language",
				Important: false,
			},
			footer: []Footer{},
		},
		{
			name: "full",
			args: args{
				message: `fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

BREAKING CHANGE: use '.use()' instea of '.load()'

before:
'''javascript
app.load({})
'''

after:
'''javascript
app.use({})
'''

Reviewed-by: Z
Refs: #123`,
			},
			want: Message{
				Header: "fix: prevent racing of requests",
				Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
				Footer: []string{
					`BREAKING CHANGE: use '.use()' instea of '.load()'

before:
'''javascript
app.load({})
'''

after:
'''javascript
app.use({})
'''`,
					`Reviewed-by: Z`,
					`Refs: #123`,
				},
			},
			header: Header{
				Type:      "fix",
				Scope:     "",
				Subject:   "prevent racing of requests",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:   "BREAKING CHANGE",
					Title: "use '.use()' instea of '.load()'",
					Content: `before:
'''javascript
app.load({})
'''

after:
'''javascript
app.use({})
'''`,
				},
				{
					Tag:     "Reviewed-by",
					Title:   "Z",
					Content: "",
				},
				{
					Tag:   "Refs",
					Title: "#123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := Parse(tt.args.message)

			assert.Equal(t, tt.want, msg)

			assert.Equal(t, tt.header, msg.GetHeader())
			assert.Equal(t, tt.footer, msg.GetFooter())
		})
	}
}
