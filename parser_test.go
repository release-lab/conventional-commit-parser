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
		want   *Message
		header Header
		footer []Footer
		Closes []string
		only   bool
	}{
		{
			name: "revert commit",
			args: args{
				message: `Revert "deprecated"
This reverts commit bf08694.`,
			},
			want: &Message{
				Header: `Revert "deprecated"`,
				Body:   "This reverts commit bf08694.",
				Footer: []string{},
			},
			header: Header{
				Type:      "revert",
				Scope:     "",
				Subject:   "deprecated",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:     "revert",
					Title:   "deprecated",
					Content: "bf08694",
				},
			},
		},
		{
			name: "revert commit with header type",
			args: args{
				message: `revert: change tittle
This reverts commit bf08694.`,
			},
			want: &Message{
				Header: `revert: change tittle`,
				Body:   "This reverts commit bf08694.",
				Footer: []string{},
			},
			header: Header{
				Type:      "revert",
				Scope:     "",
				Subject:   "change tittle",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:     "revert",
					Title:   "change tittle",
					Content: "bf08694",
				},
			},
		},
		{
			name: "revert commit with header type but no valid body",
			args: args{
				message: `revert: change tittle

revert it`,
			},
			want: &Message{
				Header: `revert: change tittle`,
				Body:   "revert it",
				Footer: []string{},
			},
			header: Header{
				Type:      "revert",
				Scope:     "",
				Subject:   "change tittle",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:     "revert",
					Title:   "change tittle",
					Content: "",
				},
			},
		},
		{
			name: "simple breaking change",
			only: true,
			args: args{
				message: `feat: rename tag:N to @N

BREAKING CHANGE: rename

'''diff
- tag:0~tag:1
+ @0~@1
'''`,
			},
			want: &Message{
				Header: `feat: rename tag:N to @N`,
				Body:   "",
				Footer: []string{
					`BREAKING CHANGE: rename

'''diff
- tag:0~tag:1
+ @0~@1
'''`,
				},
			},
			header: Header{
				Type:      "feat",
				Scope:     "",
				Subject:   "rename tag:N to @N",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:   "BREAKING CHANGE",
					Title: "rename",
					Content: `'''diff
- tag:0~tag:1
+ @0~@1
'''`,
				},
			},
		},
		{
			name: "close issue",
			only: true,
			args: args{
				message: `feat: support xxx

Closes #1, #2, #3
`,
			},
			want: &Message{
				Header: `feat: support xxx`,
				Body:   "",
				Footer: []string{"Closes #1, #2, #3"},
			},
			header: Header{
				Type:      "feat",
				Scope:     "",
				Subject:   "support xxx",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:     "Closes",
					Title:   "#1, #2, #3",
					Content: ``,
				},
			},
			Closes: []string{"#1", "#2", "#3"},
		},
		{
			name: "revert commit with other body",
			args: args{
				message: `Revert "deprecated"
This reverts commit bf08694.

Fixing these warnings, unfortunately also means, that hot code
(which reloads a shared library during runtime) can not use V
constants, because the private static C variables in the shared
library will not be initialized by _vinit(), which is only called
by the main V program.

For example in examples/hot_reload/bounce.v, using 'gx.blue',
defined as:
'    blue   = Color { r:   0, g:   0, b: 255 }'
... will instead use a const with all 0 fields (i.e. a black color).`,
			},
			want: &Message{
				Header: `Revert "deprecated"`,
				Body: `This reverts commit bf08694.

Fixing these warnings, unfortunately also means, that hot code
(which reloads a shared library during runtime) can not use V
constants, because the private static C variables in the shared
library will not be initialized by _vinit(), which is only called
by the main V program.

For example in examples/hot_reload/bounce.v, using 'gx.blue',
defined as:
'    blue   = Color { r:   0, g:   0, b: 255 }'
... will instead use a const with all 0 fields (i.e. a black color).`,
				Footer: []string{},
			},
			header: Header{
				Type:      "revert",
				Scope:     "",
				Subject:   "deprecated",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:     "revert",
					Title:   "deprecated",
					Content: "bf08694",
				},
			},
		},
		{
			name: "common commit",
			args: args{
				message: "this is a commit message",
			},
			want: &Message{
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
			want: &Message{
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
			name: "common breaking change",
			args: args{
				message: `feat(BREAKING): remove hashURL function in template render

BREAKING CHANGE:

before

'''bash
{{ hashURL .Hash}}
{{ hashURL .RevertCommitHash }}
'''

after

'''bash
{{ .HashURL }}
{{ .RevertCommitHashURL }}
'''`,
			},
			want: &Message{
				Header: "feat(BREAKING): remove hashURL function in template render",
				Body:   "",
				Footer: []string{`BREAKING CHANGE:

before

'''bash
{{ hashURL .Hash}}
{{ hashURL .RevertCommitHash }}
'''

after

'''bash
{{ .HashURL }}
{{ .RevertCommitHashURL }}
'''`},
			},
			header: Header{
				Type:      "feat",
				Scope:     "BREAKING",
				Subject:   "remove hashURL function in template render",
				Important: false,
			},
			footer: []Footer{
				{
					Tag:   "BREAKING CHANGE",
					Title: "",
					Content: `before

'''bash
{{ hashURL .Hash}}
{{ hashURL .RevertCommitHash }}
'''

after

'''bash
{{ .HashURL }}
{{ .RevertCommitHashURL }}
'''`,
				},
			},
		},
		{
			name: "Commit message with description and breaking change footer",
			args: args{
				message: `feat: allow provided config object to extend other configs

BREAKING CHANGE: 'extends' key in config file is now used for extending other config files`,
			},
			want: &Message{
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
			want: &Message{
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
			want: &Message{
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
			want: &Message{
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
			want: &Message{
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
			want: &Message{
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
			want: &Message{
				Header: "fix: prevent racing of requests",
				Body: `Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.`,
				Footer: []string{
					`BREAKING CHANGE: use '.use()' instea of '.load()'`,
					`before:
'''javascript
app.load({})
'''`,
					`after:
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
					Tag:     "BREAKING CHANGE",
					Title:   "use '.use()' instea of '.load()'",
					Content: "",
				},
				{
					Tag:   "before",
					Title: "",
					Content: `'''javascript
app.load({})
'''`,
				},
				{
					Tag:   "after",
					Title: "",
					Content: `'''javascript
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

			assert.Equal(t, tt.header, msg.ParseHeader())
			assert.Equal(t, tt.footer, msg.ParseFooter())

			if tt.Closes != nil && len(tt.Closes) != 0 {
				assert.Equal(t, tt.Closes, msg.GetCloses())
			}
		})
	}
}
